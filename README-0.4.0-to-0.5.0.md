# Changelog `set v0.4.0` to `set v0.5.0`

This document outlines the changes and justifications between the versions listed above. In general `set v0.5.0` addresses performance and optimization bottlenecks in the prior version.

There are a few breaking changes but they are mostly minor.

## Remove Unnecessary Interface & Switch to Value-Types

The `BoundMapping` interface has been redefined as a struct{}.

The `BoundMapping`, `Mapping`, and `Value` types are no longer created and returned as pointers. Where appropriate their methods use pointer receivers but do not perform `nil` checks.

These changes remove unnecessary pointer dereferences and lower the frequency of values escaping to heap.

## Altered Method Signature(s)

```go
func (m *Mapper) Bind(...) BoundMapping // v0.4.0, old

func (m Mapper) Bind(...) (BoundMapping, error) // v0.5.0, new
```

The `Bind` method now returns an error. This allows the caller to check for incorrect usage, such as supplying a read-only value, at the time of calling Bind.

Any error returned from Bind is also set internally in the BoundMapping. BoundMapping methods will return a pre-set error and no longer have to check on every call if the bound value is usable.

## Sentinal Errors

`v0.5.0` no longer depends on my `errors` package and instead uses sentinal errors. Error paths in `set` no longer produce stack traces and run more efficiently. The returned errors still contain contextual information about where and why the error occurred.

## Deprecated Behavior

The prior version allowed a `*set.Value` to be passed to `Mapper.Bind` or `BoundMapping.Rebind`. This is no longer supported.

This behavior was not publicly documented and buried inside an internal test. I doubt it was used.

However in `v0.5.0` both `Mapper.Bind` and `BoundMapping.Rebind` accept a `reflect.Value` to bind or rebind to. Therefore this behavior is still supported by passing the `set.Value.TopValue` field to bind or rebind.

Open an issue or ask me for help in the Go discord channel if this affects you.

# New Types & Subpackages

## `PreparedMapping`

`PreparedMapping` allows access to a bound value's fields similar to the `BoundMapping` type; however a PreparedMapping requires an access plan to be set with its `Plan` method. Once a plan is established the desired field names are no longer required when calling methods on the PreparedMapping.

`PreparedMapping` is more performant than `BoundMapping` when struct fields will be accessed in the same order after every call to `Rebind`.

## `coerce`

`coerce` is a new subpackage that performs type coercion. The type coercion implementation in the previous version is incredibly slow and inefficient compared to the new implementation.

I can not guarantee the new coercion logic is compatible in every way with the old coercion logic. However all of the existing tests still pass and I've added many new tests.

## `path`

`path` is a new subpackage that performs the initial mapping of Go structs. While traversing structs the `path` package builds information to traverse structs via `reflect` or `unsafe`.

This new version of `set` only traverses structs with `reflect`. However a future version may use or allow `unsafe` traversal and this functionality will be provided by `path`.

# Benchmarks

## Type Coercion

`set v0.5.0` includes a new benchmark `BenchmarkScalarTo` to evaluate the type coercion system.

`v0.4.0`

```
BenchmarkScalarTo
BenchmarkScalarTo/bool-8                  343676              3512 ns/op               0 B/op          0 allocs/op
BenchmarkScalarTo/float32-8               306393              3724 ns/op               4 B/op          1 allocs/op
BenchmarkScalarTo/float64-8               299218              3727 ns/op               8 B/op          1 allocs/op
BenchmarkScalarTo/int-8                   323304              3780 ns/op               8 B/op          1 allocs/op
BenchmarkScalarTo/int8-8                  323512              3758 ns/op               1 B/op          1 allocs/op
BenchmarkScalarTo/int16-8                 315246              3720 ns/op               2 B/op          1 allocs/op
BenchmarkScalarTo/int32-8                 307384              3689 ns/op               4 B/op          1 allocs/op
BenchmarkScalarTo/int64-8                 307186              3741 ns/op               8 B/op          1 allocs/op
BenchmarkScalarTo/uint-8                  307382              3808 ns/op               8 B/op          1 allocs/op
BenchmarkScalarTo/uint8-8                 315027              3718 ns/op               1 B/op          1 allocs/op
BenchmarkScalarTo/uint16-8                307386              3785 ns/op               2 B/op          1 allocs/op
BenchmarkScalarTo/uint32-8                323005              3730 ns/op               4 B/op          1 allocs/op
BenchmarkScalarTo/uint64-8                315001              3767 ns/op               8 B/op          1 allocs/op
BenchmarkScalarTo/string-8                159622              7247 ns/op              80 B/op         28 allocs/op
```

`v0.5.0`

```
BenchmarkScalarTo
BenchmarkScalarTo/bool-8                 2579983               461.9 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/float32-8              2274598               537.4 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/float64-8              2309199               520.8 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/int-8                  2292742               516.4 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/int8-8                 2142733               550.1 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/int16-8                2208738               542.2 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/int32-8                2229578               538.1 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/int64-8                2236003               549.8 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/uint-8                 2369776               506.1 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/uint8-8                2226345               535.8 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/uint16-8               2222169               539.6 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/uint32-8               2252427               532.2 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/uint64-8               2309792               518.5 ns/op             0 B/op          0 allocs/op
BenchmarkScalarTo/string-8               1444926               826.6 ns/op           104 B/op          8 allocs/op
```

Compared with `benchstat`

```
name                old time/op    new time/op    delta
ScalarTo/bool-8       3.60µs ± 3%    0.46µs ± 1%   -87.21%  (p=0.008 n=5+5)
ScalarTo/float32-8    3.83µs ± 2%    0.53µs ± 2%   -86.14%  (p=0.008 n=5+5)
ScalarTo/float64-8    3.74µs ± 2%    0.52µs ± 1%   -85.97%  (p=0.008 n=5+5)
ScalarTo/int-8        3.73µs ± 2%    0.51µs ± 0%   -86.31%  (p=0.016 n=5+4)
ScalarTo/int8-8       3.75µs ± 1%    0.55µs ± 1%   -85.22%  (p=0.008 n=5+5)
ScalarTo/int16-8      3.75µs ± 2%    0.55µs ± 1%   -85.42%  (p=0.008 n=5+5)
ScalarTo/int32-8      3.75µs ± 2%    0.54µs ± 1%   -85.51%  (p=0.008 n=5+5)
ScalarTo/int64-8      3.77µs ± 2%    0.54µs ± 0%   -85.79%  (p=0.008 n=5+5)
ScalarTo/uint-8       3.86µs ± 3%    0.52µs ± 1%   -86.59%  (p=0.008 n=5+5)
ScalarTo/uint8-8      3.91µs ± 3%    0.55µs ± 1%   -86.02%  (p=0.008 n=5+5)
ScalarTo/uint16-8     3.92µs ± 1%    0.54µs ± 2%   -86.13%  (p=0.008 n=5+5)
ScalarTo/uint32-8     3.85µs ± 1%    0.54µs ± 1%   -85.93%  (p=0.008 n=5+5)
ScalarTo/uint64-8     3.86µs ± 1%    0.52µs ± 0%   -86.40%  (p=0.008 n=5+5)
ScalarTo/string-8     7.31µs ± 2%    0.83µs ± 2%   -88.60%  (p=0.008 n=5+5)
```

## `BoundMapping`

`set v0.5.0` includes many new benchmarks for `BoundMapping`. This section compares these benchmarks when run with `v0.4.0` and `v0.5.0`.

`v0.4.0`

```
BenchmarkMapper/Bind_no_Rebind-8          871509              1220 ns/op             384 B/op         10 allocs/op
BenchmarkMapper/Bind_Rebind-8            1214550               998.6 ns/op           128 B/op          8 allocs/op
Benchmark_Mapper_BindPrepare/Bind_unknown1-8               30954             38601 ns/op           10772 B/op        110 allocs/op
Benchmark_Mapper_BindPrepare/Bind_unknown2-8               31384             38476 ns/op           10772 B/op        110 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Assignables_simple-8          12542539                94.97 ns/op            0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Fields_simple-8                8889804               130.4 ns/op            24 B/op          2 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Field_simple-8                 2791956               424.3 ns/op           496 B/op          4 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Set_simple-8                   2659396               448.4 ns/op           248 B/op          2 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Assignables_nested-8           5980756               201.8 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Fields_nested-8                4653588               266.1 ns/op            48 B/op          4 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Field_nested-8                 1389459               859.1 ns/op           992 B/op          8 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Set_nested-8                    878136              1330 ns/op             746 B/op          7 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Assignables_ptr_nested-8       5143212               232.7 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Fields_ptr_nested-8            4058233               294.0 ns/op            48 B/op          4 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Field_ptr_nested-8             1352497               887.2 ns/op           992 B/op          8 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Set_ptr_nested-8                878187              1366 ns/op             746 B/op          7 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Assignables_primitives-8       1931373               624.2 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Fields_primitives-8            1724206               708.6 ns/op            88 B/op         14 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Field_primitives-8              359887              3177 ns/op            3472 B/op         28 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Set_primitives-8               2138775               563.9 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Assignables_sale-8             1492143               802.8 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Fields_sale-8                  1000000              1030 ns/op             184 B/op         15 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Field_sale-8                    341296              3523 ns/op            3720 B/op         30 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Set_sale-8                      558264              2161 ns/op             840 B/op          8 allocs/op
```

`v0.5.0`

```
BenchmarkMapper/Bind_no_Rebind-8         1316692               891.6 ns/op           128 B/op          8 allocs/op
BenchmarkMapper/Bind_Rebind-8            1508492               803.9 ns/op           128 B/op          8 allocs/op
Benchmark_Mapper_BindPrepare/Bind_unknown1-8             2053950               582.2 ns/op           576 B/op          9 allocs/op
Benchmark_Mapper_BindPrepare/Bind_unknown2-8             2040548               593.8 ns/op           576 B/op          9 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Assignables_simple-8          14612013                82.91 ns/op            0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Fields_simple-8               12942205                93.27 ns/op            0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Field_simple-8                 4567026               261.0 ns/op            48 B/op          2 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Set_simple-8                   5918787               201.7 ns/op            24 B/op          1 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Assignables_nested-8           7317282               165.0 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Fields_nested-8                6391558               179.5 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Field_nested-8                 2250918               537.0 ns/op            96 B/op          4 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Set_nested-8                   2250134               530.9 ns/op            72 B/op          3 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Assignables_ptr_nested-8       6621734               176.9 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Fields_ptr_nested-8            5963674               194.0 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Field_ptr_nested-8             2128747               559.0 ns/op            96 B/op          4 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Set_ptr_nested-8               2194550               544.2 ns/op            72 B/op          3 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Assignables_primitives-8       2193958               551.2 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Fields_primitives-8            2141484               562.8 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Field_primitives-8              585465              1979 ns/op             336 B/op         14 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Set_primitives-8               2627559               457.0 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Assignables_sale-8             1664931               690.5 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Fields_sale-8                  1566954               784.1 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Field_sale-8                    512841              2309 ns/op             360 B/op         15 allocs/op
Benchmark_Mapper_BindPrepare/Bind_Set_sale-8                      878206              1347 ns/op             280 B/op          9 allocs/op
```

Compared with `benchstat`

```
name                                               old time/op    new time/op    delta
Mapper/Bind_no_Rebind-8                              1.20µs ± 2%    0.91µs ± 3%   -24.29%  (p=0.008 n=5+5)
Mapper/Bind_Rebind-8                                  986ns ± 2%     808ns ± 2%   -18.02%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_unknown1-8                  38.3µs ± 1%     0.6µs ± 0%   -98.47%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_unknown2-8                  38.3µs ± 1%     0.6µs ± 0%   -98.45%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_simple-8        93.2ns ± 1%    84.2ns ± 2%    -9.62%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Fields_simple-8              129ns ± 0%      94ns ± 2%   -27.31%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_simple-8               413ns ± 1%     261ns ± 1%   -36.74%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_simple-8                 448ns ± 1%     203ns ± 1%   -54.72%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_nested-8         198ns ± 1%     162ns ± 3%   -18.50%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Fields_nested-8              258ns ± 1%     182ns ± 2%   -29.61%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_nested-8               836ns ± 0%     526ns ± 0%   -37.11%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_nested-8                1.33µs ± 1%    0.53µs ± 0%   -60.49%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_ptr_nested-8     234ns ± 3%     173ns ± 3%   -25.93%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Fields_ptr_nested-8          296ns ± 2%     192ns ± 3%   -35.17%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_ptr_nested-8           870ns ± 1%     550ns ± 1%   -36.84%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_ptr_nested-8            1.37µs ± 1%    0.55µs ± 0%   -60.06%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_primitives-8     571ns ± 2%     546ns ± 2%    -4.47%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Fields_primitives-8          705ns ± 1%     599ns ± 2%   -14.99%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_primitives-8          3.07µs ± 0%    2.02µs ± 1%   -34.05%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_primitives-8             569ns ± 1%     454ns ± 0%   -20.13%  (p=0.016 n=5+4)
_Mapper_BindPrepare/Bind_Assignables_sale-8           765ns ± 2%     662ns ± 2%   -13.38%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Fields_sale-8               1.08µs ± 1%    0.74µs ± 5%   -31.01%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_sale-8                3.51µs ± 1%    2.34µs ± 1%   -33.46%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_sale-8                  2.09µs ± 0%    1.36µs ± 2%   -35.12%  (p=0.008 n=5+5)

name                                               old alloc/op   new alloc/op   delta
Mapper/Bind_no_Rebind-8                                384B ± 0%      128B ± 0%   -66.67%  (p=0.008 n=5+5)
Mapper/Bind_Rebind-8                                   128B ± 0%      128B ± 0%      ~     (all equal)
_Mapper_BindPrepare/Bind_unknown1-8                  10.8kB ± 0%     0.6kB ± 0%   -94.65%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_unknown2-8                  10.8kB ± 0%     0.6kB ± 0%   -94.65%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_simple-8         0.00B          0.00B           ~     (all equal)
_Mapper_BindPrepare/Bind_Fields_simple-8              24.0B ± 0%      0.0B       -100.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_simple-8                496B ± 0%       48B ± 0%   -90.32%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_simple-8                  248B ± 0%       24B ± 0%   -90.32%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_nested-8         0.00B          0.00B           ~     (all equal)
_Mapper_BindPrepare/Bind_Fields_nested-8              48.0B ± 0%      0.0B       -100.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_nested-8                992B ± 0%       96B ± 0%   -90.32%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_nested-8                  746B ± 0%       72B ± 0%   -90.35%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_ptr_nested-8     0.00B          0.00B           ~     (all equal)
_Mapper_BindPrepare/Bind_Fields_ptr_nested-8          48.0B ± 0%      0.0B       -100.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_ptr_nested-8            992B ± 0%       96B ± 0%   -90.32%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_ptr_nested-8              746B ± 0%       72B ± 0%   -90.35%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_primitives-8     0.00B          0.00B           ~     (all equal)
_Mapper_BindPrepare/Bind_Fields_primitives-8          88.0B ± 0%      0.0B       -100.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_primitives-8          3.47kB ± 0%    0.34kB ± 0%   -90.32%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_primitives-8             0.00B          0.00B           ~     (all equal)
_Mapper_BindPrepare/Bind_Assignables_sale-8           0.00B          0.00B           ~     (all equal)
_Mapper_BindPrepare/Bind_Fields_sale-8                 184B ± 0%        0B       -100.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_sale-8                3.72kB ± 0%    0.36kB ± 0%   -90.32%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_sale-8                    840B ± 0%      280B ± 0%   -66.67%  (p=0.008 n=5+5)

name                                               old allocs/op  new allocs/op  delta
Mapper/Bind_no_Rebind-8                                10.0 ± 0%       8.0 ± 0%   -20.00%  (p=0.008 n=5+5)
Mapper/Bind_Rebind-8                                   8.00 ± 0%      8.00 ± 0%      ~     (all equal)
_Mapper_BindPrepare/Bind_unknown1-8                     110 ± 0%         9 ± 0%   -91.82%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_unknown2-8                     110 ± 0%         9 ± 0%   -91.82%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_simple-8          0.00           0.00           ~     (all equal)
_Mapper_BindPrepare/Bind_Fields_simple-8               2.00 ± 0%      0.00       -100.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_simple-8                4.00 ± 0%      2.00 ± 0%   -50.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_simple-8                  2.00 ± 0%      1.00 ± 0%   -50.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_nested-8          0.00           0.00           ~     (all equal)
_Mapper_BindPrepare/Bind_Fields_nested-8               4.00 ± 0%      0.00       -100.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_nested-8                8.00 ± 0%      4.00 ± 0%   -50.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_nested-8                  7.00 ± 0%      3.00 ± 0%   -57.14%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_ptr_nested-8      0.00           0.00           ~     (all equal)
_Mapper_BindPrepare/Bind_Fields_ptr_nested-8           4.00 ± 0%      0.00       -100.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_ptr_nested-8            8.00 ± 0%      4.00 ± 0%   -50.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_ptr_nested-8              7.00 ± 0%      3.00 ± 0%   -57.14%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Assignables_primitives-8      0.00           0.00           ~     (all equal)
_Mapper_BindPrepare/Bind_Fields_primitives-8           14.0 ± 0%       0.0       -100.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_primitives-8            28.0 ± 0%      14.0 ± 0%   -50.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_primitives-8              0.00           0.00           ~     (all equal)
_Mapper_BindPrepare/Bind_Assignables_sale-8            0.00           0.00           ~     (all equal)
_Mapper_BindPrepare/Bind_Fields_sale-8                 15.0 ± 0%       0.0       -100.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Field_sale-8                  30.0 ± 0%      15.0 ± 0%   -50.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Bind_Set_sale-8                    8.00 ± 0%      9.00 ± 0%   +12.50%  (p=0.008 n=5+5)
```

## `PreparedMapping`

The previous section highlights the improved performance of the BoundMapping type. This section compares the new `PreparedMapping` type with the improved BoundMapping implementation.

`PreparedMapping`

```
BenchmarkMapper/Prepare_Rebind-8         1918274               615.4 ns/op           128 B/op          8 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_unknown1-8          4470338               267.8 ns/op           256 B/op          4 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_unknown2-8          4394802               271.8 ns/op           256 B/op          4 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Assignables_simple-8               18621049                63.62 ns/op            0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Fields_simple-8                    16888491                71.54 ns/op            0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Field_simple-8                      4905577               243.6 ns/op            48 B/op          2 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Set_simple-8                        6336757               185.9 ns/op            24 B/op          1 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Assignables_nested-8                9975766               121.7 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Fields_nested-8                     8327191               141.4 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Field_nested-8                      2453205               488.3 ns/op            96 B/op          4 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Set_nested-8                        2398924               495.6 ns/op            72 B/op          3 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Assignables_ptr_nested-8            8590812               138.5 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Fields_ptr_nested-8                 7554633               152.3 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Field_ptr_nested-8                  2311070               519.9 ns/op            96 B/op          4 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Set_ptr_nested-8                    2350860               517.4 ns/op            72 B/op          3 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Assignables_primitives-8            3920596               307.2 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Fields_primitives-8                 3121963               398.8 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Field_primitives-8                   641995              1828 ns/op             336 B/op         14 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Set_primitives-8                    3842169               302.5 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Assignables_sale-8                  2836902               427.5 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Fields_sale-8                       2350030               521.0 ns/op             0 B/op          0 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Field_sale-8                         585480              2038 ns/op             360 B/op         15 allocs/op
Benchmark_Mapper_BindPrepare/Prepare_Set_sale-8                          1000000              1124 ns/op             280 B/op          9 allocs/op
```

Compared against BoundMapping with `benchstat`

```
name                                          old time/op    new time/op    delta
Mapper/Rebind-8                                  808ns ± 2%     614ns ± 1%  -24.05%  (p=0.008 n=5+5)
_Mapper_BindPrepare/unaddr1-8                    331ns ± 0%     232ns ± 0%  -29.89%  (p=0.008 n=5+5)
_Mapper_BindPrepare/unaddr2-8                    335ns ± 1%     233ns ± 0%  -30.58%  (p=0.008 n=5+5)
_Mapper_BindPrepare/unknown1-8                   584ns ± 0%     269ns ± 1%  -54.01%  (p=0.008 n=5+5)
_Mapper_BindPrepare/unknown2-8                   593ns ± 0%     269ns ± 3%  -54.56%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Assignables_simple-8        84.2ns ± 2%    63.4ns ± 2%  -24.71%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Fields_simple-8             93.5ns ± 2%    71.8ns ± 1%  -23.28%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Field_simple-8               261ns ± 1%     245ns ± 0%   -6.32%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Set_simple-8                 203ns ± 1%     188ns ± 1%   -7.30%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Assignables_nested-8         162ns ± 3%     121ns ± 1%  -24.89%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Fields_nested-8              182ns ± 2%     142ns ± 2%  -21.56%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Field_nested-8               526ns ± 0%     488ns ± 1%   -7.23%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Set_nested-8                 525ns ± 0%     501ns ± 2%   -4.62%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Assignables_ptr_nested-8     173ns ± 3%     139ns ± 1%  -19.71%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Fields_ptr_nested-8          192ns ± 3%     152ns ± 1%  -21.04%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Field_ptr_nested-8           550ns ± 1%     511ns ± 0%   -6.97%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Set_ptr_nested-8             546ns ± 0%     517ns ± 0%   -5.28%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Assignables_primitives-8     546ns ± 2%     304ns ± 1%  -44.27%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Fields_primitives-8          599ns ± 2%     396ns ± 2%  -34.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Field_primitives-8          2.02µs ± 1%    1.76µs ± 0%  -12.96%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Set_primitives-8             454ns ± 0%     303ns ± 3%  -33.41%  (p=0.016 n=4+5)
_Mapper_BindPrepare/Assignables_sale-8           662ns ± 2%     421ns ± 1%  -36.49%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Fields_sale-8                745ns ± 5%     514ns ± 2%  -31.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Field_sale-8                2.34µs ± 1%    1.99µs ± 1%  -15.00%  (p=0.008 n=5+5)
_Mapper_BindPrepare/Set_sale-8                  1.36µs ± 2%    1.13µs ± 2%  -16.80%  (p=0.008 n=5+5)
```
