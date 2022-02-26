# Log libs benchmarks

## v0.2.0

> record ad 2022.02.26

```text
$ go test -v -cpu=4 -run=none -bench=. -benchtime=10s -benchmem bench_loglibs_test.go
goos: windows
goarch: amd64                               
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
BenchmarkZapNegative
BenchmarkZapNegative-4                  139243226               86.39 ns/op          192 B/op          1 allocs/op
BenchmarkZeroLogNegative
BenchmarkZeroLogNegative-4              1000000000               8.302 ns/op           0 B/op          0 allocs/op
BenchmarkPhusLogNegative
BenchmarkPhusLogNegative-4              1000000000               8.989 ns/op           0 B/op          0 allocs/op
BenchmarkLogrusNegative
BenchmarkGookitSlogNegative-4           38300540               323.3 ns/op           221 B/op          5 allocs/op
BenchmarkZapPositive
BenchmarkZapPositive-4                  14453001               828.1 ns/op           192 B/op          1 allocs/op
BenchmarkZeroLogPositive
BenchmarkZeroLogPositive-4              28671724               420.9 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogPositive
BenchmarkPhusLogPositive-4              45619569               261.9 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusPositive
BenchmarkLogrusPositive-4                5092164              2366 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogPositive
BenchmarkGookitSlogPositive-4            3184557              3754 ns/op             856 B/op         13 allocs/op
PASS
ok      command-line-arguments  135.460s
```

## v0.1.5

> record ad 2022.02.26

```text
$ go test -v -cpu=4 -run=none -bench=. -benchtime=10s -benchmem bench_loglibs_test.go
goos: windows
goarch: amd64
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
BenchmarkZapNegative
BenchmarkZapNegative-4                  137676860               86.43 ns/op          192 B/op          1 allocs/op
BenchmarkZeroLogNegative
BenchmarkZeroLogNegative-4              1000000000               8.284 ns/op           0 B/op          0 allocs/op
BenchmarkPhusLogNegative
BenchmarkZapPositive-4                  14250313               831.7 ns/op           192 B/op          1 allocs/op
BenchmarkZeroLogPositive
BenchmarkZeroLogPositive-4              28183436               426.0 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogPositive
BenchmarkPhusLogPositive-4              44034984               258.7 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusPositive
BenchmarkLogrusPositive-4                5005593              2421 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogPositive
BenchmarkGookitSlogPositive-4            1714084              7029 ns/op            4480 B/op         45 allocs/op
PASS
ok      command-line-arguments  138.199s
```