go-and
======

[![Go Reference](https://pkg.go.dev/badge/github.com/bwesterb/go-and.svg)](https://pkg.go.dev/github.com/bwesterb/go-and)

Fast bitwise and, or, xor, andn, popcount and memset for `[]byte` slices.

```go
import "github.com/bwesterb/go-and"

func main() {
    var a, b, dst []byte

    // ...

    // Sets dst to the bitwise and of a and b
    and.And(dst, a, b)
}
```

Makes use of AVX2 on AMD64 and NEON on ARM64.

## Benchmarks

Created using `./benchmark.sh`.

This shows three benchmarks:

* `naive` is a simple for loop doing one byte at a time.
* `purego` are our slightly optimized versions that work on uint64s instead of bytes.
* `asm` are the AVX2 implementations and the reason to use this library.

### Intel Core i9-13900
```
goos: linux
goarch: amd64
pkg: github.com/bwesterb/go-and
cpu: 13th Gen Intel(R) Core(TM) i9-13900
             │    naive     │                purego                │                 asm                 │
             │    sec/op    │    sec/op     vs base                │   sec/op     vs base                │
And-32         6844.5n ± 1%   2837.5n ± 0%  -58.54% (p=0.000 n=10)   365.9n ± 1%  -94.65% (p=0.000 n=10)
Or-32          7098.5n ± 2%   2837.0n ± 0%  -60.03% (p=0.000 n=10)   371.7n ± 2%  -94.76% (p=0.000 n=10)
Xor-32         6875.5n ± 2%   2838.0n ± 1%  -58.72% (p=0.000 n=10)   368.9n ± 3%  -94.64% (p=0.000 n=10)
AndNot-32      7587.0n ± 0%   3100.5n ± 0%  -59.13% (p=0.000 n=10)   367.3n ± 3%  -95.16% (p=0.000 n=10)
AnyMasked-32   6337.5n ± 4%   1604.0n ± 2%  -74.69% (p=0.000 n=10)   369.3n ± 1%  -94.17% (p=0.000 n=10)
Any-32         5673.0n ± 0%   1425.5n ± 0%  -74.87% (p=0.000 n=10)   208.8n ± 0%  -96.32% (p=0.000 n=10)
geomean         6.708µ         2.335µ       -65.20%                  335.3n       -95.00%

             │    naive     │                 purego                 │                   asm                    │
             │     B/s      │      B/s       vs base                 │      B/s        vs base                  │
And-32         4.354Gi ± 1%   10.504Gi ± 0%  +141.23% (p=0.000 n=10)    81.456Gi ± 1%  +1770.70% (p=0.000 n=10)
Or-32          4.198Gi ± 2%   10.507Gi ± 0%  +150.25% (p=0.000 n=10)    80.188Gi ± 2%  +1809.97% (p=0.000 n=10)
Xor-32         4.335Gi ± 2%   10.502Gi ± 1%  +142.27% (p=0.000 n=10)    80.798Gi ± 2%  +1764.01% (p=0.000 n=10)
AndNot-32      3.928Gi ± 0%    9.612Gi ± 0%  +144.71% (p=0.000 n=10)    81.147Gi ± 3%  +1965.81% (p=0.000 n=10)
AnyMasked-32   4.703Gi ± 4%   18.582Gi ± 2%  +295.16% (p=0.000 n=10)    80.698Gi ± 1%  +1616.04% (p=0.000 n=10)
Any-32         5.254Gi ± 0%   20.906Gi ± 0%  +297.94% (p=0.000 n=10)   142.718Gi ± 0%  +2616.60% (p=0.000 n=10)
geomean        4.443Gi         12.77Gi       +187.34%                    88.89Gi       +1900.72%
```

### Apple M2 Pro

```
goos: darwin
goarch: arm64
pkg: github.com/bwesterb/go-and
          │    naive     │                purego                │                 asm                 │
          │    sec/op    │    sec/op     vs base                │   sec/op     vs base                │
And-12      9891.5n ± 0%   2248.0n ± 1%  -77.27% (p=0.000 n=10)   511.2n ± 1%  -94.83% (p=0.000 n=10)
Or-12       9824.5n ± 1%   2267.5n ± 1%  -76.92% (p=0.000 n=10)   511.2n ± 1%  -94.80% (p=0.000 n=10)
Xor-12      9826.5n ± 0%   2232.0n ± 1%  -77.29% (p=0.000 n=10)   509.8n ± 1%  -94.81% (p=0.000 n=10)
AndNot-12   9928.5n ± 1%   2251.0n ± 1%  -77.33% (p=0.000 n=10)   581.7n ± 0%  -94.14% (p=0.000 n=10)
Memset-12   294.96µ ± 1%    79.18µ ± 2%  -73.15% (p=0.000 n=10)   10.83µ ± 7%  -96.33% (p=0.000 n=10)
Popcnt-12   288.90µ ± 1%    45.24µ ± 0%  -84.34% (p=0.000 n=10)   11.30µ ± 0%  -96.09% (p=0.000 n=10)
geomean      30.52µ         6.716µ       -77.99%                  1.455µ       -95.23%

          │    naive     │                 purego                 │                   asm                   │
          │     B/s      │      B/s       vs base                 │      B/s       vs base                  │
And-12      3.013Gi ± 0%   13.256Gi ± 1%  +339.97% (p=0.000 n=10)   58.297Gi ± 1%  +1834.87% (p=0.000 n=10)
Or-12       3.033Gi ± 1%   13.144Gi ± 1%  +333.30% (p=0.000 n=10)   58.303Gi ± 1%  +1821.98% (p=0.000 n=10)
Xor-12      3.033Gi ± 0%   13.353Gi ± 1%  +340.28% (p=0.000 n=10)   58.462Gi ± 1%  +1827.64% (p=0.000 n=10)
AndNot-12   3.002Gi ± 1%   13.240Gi ± 1%  +341.07% (p=0.000 n=10)   51.233Gi ± 0%  +1606.77% (p=0.000 n=10)
Memset-12   3.157Gi ± 1%   11.761Gi ± 2%  +272.50% (p=0.000 n=10)   86.041Gi ± 7%  +2625.05% (p=0.000 n=10)
Popcnt-12   3.224Gi ± 1%   20.585Gi ± 0%  +538.57% (p=0.000 n=10)   82.407Gi ± 0%  +2456.30% (p=0.000 n=10)
geomean     3.076Gi         13.98Gi       +354.43%                   64.53Gi       +1997.80%
```
