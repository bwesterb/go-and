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
          │    naive     │               purego                │                 asm                 │
          │    sec/op    │   sec/op     vs base                │   sec/op     vs base                │
And-32      264.51µ ± 6%   64.69µ ± 1%  -75.54% (p=0.000 n=10)   24.43µ ± 3%  -90.77% (p=0.000 n=10)
Or-32       274.06µ ± 5%   64.89µ ± 2%  -76.32% (p=0.000 n=10)   24.30µ ± 1%  -91.13% (p=0.000 n=10)
AndNot-32   309.01µ ± 0%   73.10µ ± 1%  -76.34% (p=0.000 n=10)   24.52µ ± 2%  -92.07% (p=0.000 n=10)
Memset-32   225.74µ ± 4%   56.64µ ± 1%  -74.91% (p=0.000 n=10)   15.77µ ± 1%  -93.01% (p=0.000 n=10)
Popcnt-32   128.45µ ± 1%   69.35µ ± 0%  -46.01% (p=0.000 n=10)   31.80µ ± 3%  -75.24% (p=0.000 n=10)
geomean      230.4µ        65.50µ       -71.58%                  23.59µ       -89.76%

          │    naive     │                 purego                 │                   asm                   │
          │     B/s      │      B/s       vs base                 │      B/s       vs base                  │
And-32      3.521Gi ± 6%   14.397Gi ± 2%  +308.89% (p=0.000 n=10)   38.129Gi ± 3%   +982.90% (p=0.000 n=10)
Or-32       3.398Gi ± 6%   14.353Gi ± 2%  +322.36% (p=0.000 n=10)   38.319Gi ± 1%  +1027.60% (p=0.000 n=10)
AndNot-32   3.014Gi ± 0%   12.740Gi ± 1%  +322.71% (p=0.000 n=10)   37.988Gi ± 2%  +1160.45% (p=0.000 n=10)
Memset-32   4.126Gi ± 3%   16.444Gi ± 1%  +298.59% (p=0.000 n=10)   59.051Gi ± 1%  +1331.33% (p=0.000 n=10)
Popcnt-32   7.251Gi ± 1%   13.428Gi ± 0%   +85.20% (p=0.000 n=10)   29.288Gi ± 3%   +303.94% (p=0.000 n=10)
geomean     4.042Gi         14.22Gi       +251.80%                   39.49Gi        +876.93%
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
