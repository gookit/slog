# Log libs benchmarks

Run benchmark: `make test-bench`

> **Note**: on each test will update all package to latest.

## v0.5.1 - 2023.04.13

> **Note**: test and record ad 2023.04.13

```shell
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-3740QM CPU @ 2.70GHz
BenchmarkZapNegative
BenchmarkZapNegative-4                   8381674              1429 ns/op             216 B/op          3 allocs/op
BenchmarkZapSugarNegative
BenchmarkZapSugarNegative-4              8655980              1383 ns/op             104 B/op          4 allocs/op
BenchmarkZeroLogNegative
BenchmarkZeroLogNegative-4              14173719               849.8 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogNegative
BenchmarkPhusLogNegative-4              27456256               451.2 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusNegative
BenchmarkLogrusNegative-4                2550771              4784 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogNegative
BenchmarkGookitSlogNegative-4            8798220              1375 ns/op             120 B/op          3 allocs/op
BenchmarkZapPositive
BenchmarkZapPositive-4                  10302483              1167 ns/op             192 B/op          1 allocs/op
BenchmarkZapSugarPositive
BenchmarkZapSugarPositive-4              3833311              3154 ns/op             344 B/op          7 allocs/op
BenchmarkZeroLogPositive
BenchmarkZeroLogPositive-4              14120524               846.7 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogPositive
BenchmarkPhusLogPositive-4              27152686               434.9 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusPositive
BenchmarkLogrusPositive-4                2601892              4691 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogPositive
BenchmarkGookitSlogPositive-4            8997104              1340 ns/op             120 B/op          3 allocs/op
PASS
ok      command-line-arguments  167.095s
```

## v0.3.5 - 2022.11.08

> **Note**: test and record ad 2022.11.08

```shell
% make test-bench
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-3740QM CPU @ 2.70GHz
BenchmarkZapNegative
BenchmarkZapNegative-4                  123297997              110.4 ns/op           192 B/op          1 allocs/op
BenchmarkZeroLogNegative
BenchmarkZeroLogNegative-4              891508806               13.36 ns/op            0 B/op          0 allocs/op
BenchmarkPhusLogNegative
BenchmarkPhusLogNegative-4              811990076               14.74 ns/op            0 B/op          0 allocs/op
BenchmarkLogrusNegative
BenchmarkLogrusNegative-4               242633541               49.40 ns/op           16 B/op          1 allocs/op
BenchmarkGookitSlogNegative
BenchmarkGookitSlogNegative-4           29102253               422.8 ns/op           125 B/op          4 allocs/op
BenchmarkZapPositive
BenchmarkZapPositive-4                   9772791              1194 ns/op             192 B/op          1 allocs/op
BenchmarkZeroLogPositive
BenchmarkZeroLogPositive-4              13944360               856.8 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogPositive
BenchmarkPhusLogPositive-4              27839614               431.2 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusPositive
BenchmarkLogrusPositive-4                2621076              4583 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogPositive
BenchmarkGookitSlogPositive-4            8908768              1359 ns/op             149 B/op          5 allocs/op
PASS
ok      command-line-arguments  149.379s
```

## v0.3.0

> **Note**: test and record ad 2022.04.27

```shell
% make test-bench
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-3740QM CPU @ 2.70GHz
BenchmarkZapNegative
BenchmarkZapNegative-4                  128133166               93.97 ns/op          192 B/op          1 allocs/op
BenchmarkZeroLogNegative
BenchmarkZeroLogNegative-4              909583207               13.41 ns/op            0 B/op          0 allocs/op
BenchmarkPhusLogNegative
BenchmarkPhusLogNegative-4              784099310               15.24 ns/op            0 B/op          0 allocs/op
BenchmarkLogrusNegative
BenchmarkLogrusNegative-4               289939296               41.60 ns/op           16 B/op          1 allocs/op
BenchmarkGookitSlogNegative
BenchmarkGookitSlogNegative-4           29131203               417.4 ns/op           125 B/op          4 allocs/op
BenchmarkZapPositive
BenchmarkZapPositive-4                   9910075              1219 ns/op             192 B/op          1 allocs/op
BenchmarkZeroLogPositive
BenchmarkZeroLogPositive-4              13966810               871.0 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogPositive
BenchmarkPhusLogPositive-4              26743148               446.2 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusPositive
BenchmarkLogrusPositive-4                2658482              4481 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogPositive
BenchmarkGookitSlogPositive-4            8349562              1441 ns/op             165 B/op          6 allocs/op
PASS
ok      command-line-arguments  146.669s
```

### beta 2022.04.17

> **Note**: test and record ad 2022.04.17

```shell
$ go test -v -cpu=4 -run=none -bench=. -benchtime=10s -benchmem bench_loglibs_test.go
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-3740QM CPU @ 2.70GHz
BenchmarkZapNegative
BenchmarkZapNegative-4                  130808992               91.91 ns/op          192 B/op          1 allocs/op
BenchmarkZeroLogNegative
BenchmarkZeroLogNegative-4              914445844               13.19 ns/op            0 B/op          0 allocs/op
BenchmarkPhusLogNegative
BenchmarkPhusLogNegative-4              792539167               15.32 ns/op            0 B/op          0 allocs/op
BenchmarkLogrusNegative
BenchmarkLogrusNegative-4               289393606               40.61 ns/op           16 B/op          1 allocs/op
BenchmarkGookitSlogNegative
BenchmarkGookitSlogNegative-4           29522170               405.3 ns/op           125 B/op          4 allocs/op
BenchmarkZapPositive
BenchmarkZapPositive-4                   9113048              1283 ns/op             192 B/op          1 allocs/op
BenchmarkZeroLogPositive
BenchmarkZeroLogPositive-4              14691699               797.0 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogPositive
BenchmarkPhusLogPositive-4              27634338               424.5 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusPositive
BenchmarkLogrusPositive-4                2734669              4363 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogPositive
BenchmarkGookitSlogPositive-4            7740348              1563 ns/op             165 B/op          6 allocs/op
PASS
ok      command-line-arguments  145.175s

```

## v0.2.1

> **Note**: test and record ad 2022.04.17

```shell
$ go test -v -cpu=4 -run=none -bench=. -benchtime=10s -benchmem bench_loglibs_test.go
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-3740QM CPU @ 2.70GHz
BenchmarkZapNegative
BenchmarkZapNegative-4                  125500471              125.8 ns/op           192 B/op          1 allocs/op
BenchmarkZeroLogNegative
BenchmarkZeroLogNegative-4              839046109               13.71 ns/op            0 B/op          0 allocs/op
BenchmarkPhusLogNegative
BenchmarkPhusLogNegative-4              757766400               15.56 ns/op            0 B/op          0 allocs/op
BenchmarkLogrusNegative
BenchmarkLogrusNegative-4               253178256               47.12 ns/op           16 B/op          1 allocs/op
BenchmarkGookitSlogNegative
BenchmarkGookitSlogNegative-4           30091606               401.9 ns/op            45 B/op          3 allocs/op
BenchmarkZapPositive
BenchmarkZapPositive-4                   9761935              1216 ns/op             192 B/op          1 allocs/op
BenchmarkZeroLogPositive
BenchmarkZeroLogPositive-4              13860344               837.1 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogPositive
BenchmarkPhusLogPositive-4              27666529               447.8 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusPositive
BenchmarkLogrusPositive-4                2705653              4403 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogPositive
BenchmarkGookitSlogPositive-4            1836384              6882 ns/op             680 B/op         11 allocs/op
PASS
ok      command-line-arguments  156.038s
```

## v0.2.0

> record ad 2022.02.26

```shell
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

```shell
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