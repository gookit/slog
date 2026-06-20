# gookit/slog 代码分析报告

> 分析日期:2026-06-20
> 范围:`logger` / `record` / `formatter` / `handler` / `processor` / `rotatefile` / `bufwrite`
> 方法:核心代码通读 + 对 3 个关键问题编写验证程序实测确认(含 `go test -race`)

## 总体评价

API 设计丰富、扩展性好,Handler / Processor / Formatter 三层抽象清晰。
主要短板集中在 **并发安全** 和 **热路径性能** 两方面,存在数个真实缺陷。

---

## 🔴 严重问题(已实测确认的 Bug)

### P0-1. Formatter 把已归还 pool 的 buffer 返回给调用方 → 多 logger 并发数据竞争

文件:`formatter_text.go:148-214`、`formatter_json.go:121-135`

```go
func (f *TextFormatter) Format(r *Record) ([]byte, error) {
    buf := textPool.Get()
    defer textPool.Put(buf)   // 函数返回时就把 buffer 还回池
    ...
    return buf.B, nil          // 却把 buffer 底层数组返回给了调用方
}
```

`textPool` / `jsonPool` 是**包级全局**变量。每个 logger 的 `l.mu` 只能串行化自身,
但不同 logger 之间共享同一个池。logger A 的 `Format` 返回 `buf.B` 并归还池后,
handler 仍在 `Output.Write(bts)` 读取它时,logger B 的 `Format` 从池里拿到同一 buf
并改写 → 读写竞争 + 日志内容错乱。

**实测**:两个各自正确加锁的 SugaredLogger 并发写,`go test -race` 报 `DATA RACE`:
- Write: `TextFormatter.Format() formatter_text.go:170`
- Previous read: handler 仍在读返回的 `bts`

**修复方向**:不要 `defer Put` 后返回 `buf.B`。让 handler 持有 buffer 生命周期,
或最简方案 `return append([]byte(nil), buf.B...)` 返回独立拷贝。

### P0-2. `GlobalFields` 被引用共享进每条 record → 被 processor / AddField 污染、跨请求泄漏

文件:`logger.go:98-104`、`record.go:279-286`

```go
func (l *Logger) newRecord() *Record {
    r := l.recordPool.Get().(*Record)
    r.Fields = l.GlobalFields   // 直接引用同一张 map,不是拷贝
    return r
}
// record.go
func (r *Record) AddField(name string, val any) *Record {
    if r.Fields == nil { r.Fields = make(M, 8) }
    r.Fields[name] = val        // 直接就地写 → 写进共享的 GlobalFields
    return r
}
```

任何 processor 调 `AddField/AddFields`(内置 `AddHostname`、`AddUniqueID`、
`AppendCtxKeys`)都会写进共享的 GlobalFields。

**实测**:设 `GlobalFields={app:demo}`,加一个 `AddField("perRecord",...)` 的 processor,
打一条日志后 `GlobalFields` 变成 `{app:demo, perRecord:...}`,被永久污染。
用 `AppendCtxKeys` 会把请求级 ctx 值泄漏到之后所有日志,且跨 goroutine 数据竞争。

**修复方向**:`newRecord` 不直接赋引用;在格式化阶段单独只读合并 GlobalFields,
或浅拷贝后再赋给 record,使 `AddField` 不会污染全局。

### P0-3. rotatefile 异步清理用读锁(RLock)做初始化 + 无锁置 nil

文件:`rotatefile/writer.go:324-367`、`writer.go:147`、`writer.go:361`

`asyncClean` 用 `d.mu.RLock()`(共享锁)保护 `cleanCh/stopCh` 的懒初始化与起 goroutine
——RLock 不提供互斥,两个并发 `Write` 可同时通过 nil 检查、重复建 channel/起 goroutine。
配套地,消费 goroutine `d.cleanCh = nil`(361)与 `close()` 中 `d.stopCh = nil`(147)
都是无锁写,与 `Write → asyncClean` 的读竞争。

**修复方向**:初始化用 `sync.Once` 或写锁 `Lock()`;`cleanCh/stopCh` 读写统一进锁或用 atomic。

---

## 🟠 设计不合理

### P1-4. 没有 logger 级的"级别快速门" → 被过滤的日志仍付出完整格式化代价

文件:`record.go:348-371`、`logger_write.go:60-95`

级别判断只在 `writeRecord` 内逐 handler 做,而消息格式化在更早的 `record.log/logf`
就执行了:

```go
func (r *Record) log(level Level, args []any) {
    r.Message = formatArgsWithSpaces(args) // 总是先格式化
    r.logger.writeRecord(level, r)         // 进来才判级别 + 还要抢 l.mu
}
```

**实测**:level=ErrorLevel,循环 100 次 `Debug(...)`,输出 0 字节,
但参数里的 `Stringer.String()` 被调用了 100 次。即生产环境每条被关掉的 Debug 仍会
pool Get → `fmt.Sprintf`/反射 → 抢全局锁 → pool Put。

**修复方向**:给 `Logger` 维护聚合的最大可处理级别(注册 handler 时更新),
在 `log/logf` 入口先判级别立即 return,绕过格式化与加锁。

### P1-5. 整个格式化 + 写入都在单把 `l.mu` 下串行

`writeRecord` 持锁期间做了 caller 栈回溯、processors、格式化(CPU 密集)和 IO。
高并发下吞吐被这把锁卡死。格式化本可在锁外完成(每条 record 独立),
仅对最终 `Write` 串行/或交给 handler 自己的锁。

### P1-6. `WithField/WithFields/...` 从池取一个 record 又立刻丢弃

文件:`logger.go:439-452`、`record.go:156-186`

`l.WithField` = `l.newRecord().WithField(...)`,而 `WithFields` 内部 `r.Copy()`
新建堆对象返回——从池取出的那个 record 既没用也没归还
(`// defer l.releaseRecord(r)` 被注释掉),且 `Copy()` 每次深拷贝 Data/Fields/Extra
三张 map。链式 `WithField().Info()` 这种常见写法白白做一次 pool churn + 三次 map 拷贝。

---

## 🟡 性能优化点

### P2-7. `EncodeToString` 对 `M` 多一层 `SafeString` 间接

文件:`util.go:124-129`

`EncodeToString` 用 `v.(map[string]any)` 断言,但 `Record.Data/Extra` 都是命名类型
`M`(`type M map[string]any`),该断言为 `false`。

> 订正:`mapToString` **并非死代码** —— `M` 实现了 `String()`,实际路径是
> `EncodeToString(M)` → `SafeString(M)` → `M.String()` → `mapToString`,输出正确。
> 仅多一层 Stringer 分发的开销。

修法:断言里加 `case M:` 直接走 `mapToString`(输出完全一致,省去间接)。

### P2-8. JSONFormatter 先建 `map[string]any` 再编码

文件:`formatter_json.go:80-135`(代码内已有 `// TODO perf` 注释)
额外一次 map 分配 + 反射编码,可改为直接往 buffer 拼 JSON。

### P2-9. 默认 `ReportCaller=true`

文件:`logger.go:84`、`util.go:44`
每条实际处理的日志都 `runtime.Callers` 栈回溯(有 alloc)。多数库默认关闭;
可考虑默认 false 或在文档强提示成本。

---

## 🟢 小问题 / 代码质量

- `bufwrite/line_writer.go:141-158`:大写入分支返回值可能 **> len(p)**,违反 `io.Writer`
  契约。另 `Reset()`(`line_writer.go:63`)把 `buf = buf[:0]`,之后 `Write` 的
  `copy(b.buf[b.n:], p)` 拷进零长 slice → 复位后写入静默丢失。
- `logger.go:356-360` `LastErr` 无锁读+清空 `l.err`;`Close()` 读写 `l.closed` 也无锁。
- API 冗余:`New` / `NewWithConfig` 完全相同(`logger.go:62/72`);
  `AddHandler/PushHandler/...`、`WithCtx/WithContext` 等大量同义 alias。
- `record.go` 有大段被注释的 TODO 代码(`Obj/Any/Str/Int` builder 等),建议清理或落地。
- `rotatefile/writer.go:528` `buildFilePath` 硬编码 `"%s/%s"`,跨平台建议 `filepath.Join`。

---

## 修复优先级

| 优先级 | 问题 | 影响 |
|---|---|---|
| P0 | #1 Formatter 返回已归还 buffer | 多 logger 并发必现数据竞争/日志错乱 |
| P0 | #2 GlobalFields 共享污染 | 数据泄漏 + 竞争 |
| P0 | #3 rotatefile 清理竞态 | 文件轮转并发竞争 |
| P1 | #4 级别门前置 | 生产热路径性能 |
| P1 | #5 锁外格式化 / #6 WithXxx pool churn | 吞吐与分配优化 |
| P2 | #7 死代码断言 / #8 JSON 拼接 / #9 caller 默认 | 分配与 CPU |
| P3 | 🟢 各项代码质量 | 健壮性 / 可维护性 |

---

## 补充:更深层的既有竞态(超出已修范围,待决策)

以下为 `go test -race ./...` 暴露、且 **master 上即存在** 的竞态(CI 用 `go test ./...`
不带 -race,故未覆盖):

- **rotatefile + CloseLock(`TestIssues_121`)**:作为 slog handler 使用时 `CloseLock=true`,
  写入/轮转依赖**外部 logger 锁**串行化,但异步清理 goroutine 运行在该锁之外,
  访问 `d.file/d.path/d.written` 等字段 → 竞态。彻底修复需让清理协程与外部锁协调
  (或 rotatefile 始终用自身锁),属较大重构。
- **Record 跨 goroutine 并发复用(`TestRecord_useMultiTimes`)**:Record 本就非线程安全;
  测试对单个 Record 并发写。应在文档明确「Record 不可跨 goroutine 共享」,或测试改造。

## 待决策项(影响面大/改变默认行为,建议确认后再做)

- **P1-5 锁外格式化**:把格式化移出 `writeRecord` 的大锁,涉及核心写路径重构(>100 行),风险较高。
- **P1-6 WithXxx 的 Copy 优化**:`Copy()` 对空 Data/Extra 也分配 map;改为空则留 nil 可省分配,
  但会使 JSON 空值从 `{}` 变 `null`(行为变更),需确认。
- **P2-8 JSONFormatter 直接拼 buffer**:中等改造,收益有限。
- **P2-9 `ReportCaller` 默认值**:默认 `true` 导致每条实际处理日志做栈回溯;改默认 `false` 是行为变更。

## 处理进度

- [x] P0-1 Formatter buffer 生命周期:Text/JSON Format 返回独立拷贝 + `-race` 回归测试
- [x] P0-2 GlobalFields 共享污染:newRecord 浅拷贝全局字段(空则保持 nil)+ 回归测试
- [x] P0-3 rotatefile asyncClean 竞态:sync.Once+写锁初始化、Close 时 join 清理 goroutine(cleanWg)、doClean 读 `d.path` 加读锁;附带修复 `FilesClear` daemon 的 `quitDaemon` 竞态、`MockClocker` 线程安全
  - 注:剩余 `go test -race ./rotatefile/` 偶发失败源于**测试代码**在 writer 存活期并发改写共享 `*Config`(配置应在 Create 后只读),属测试重构跟进项;CI(`go test ./...`)不带 -race,常规测试全绿
- [x] P1-4 级别快速门前置:`Logger.shouldHandle` 入口门控(禁用级别 172→37 ns/op、4→2 allocs);Panic/Fatal 始终放行
- [ ] P1-5 锁外格式化
- [ ] P1-6 WithXxx pool 复用
- [x] P2-7 EncodeToString 加 `case M`(订正:非死代码,经 `M.String()` 可达;此为省间接的等价优化)
- [ ] P2-8 JSONFormatter 直接拼接
- [ ] P2-9 ReportCaller 默认值评估
- [ ] P3 代码质量清理
