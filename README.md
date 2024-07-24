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
          │    naive     │               purego                │                 asm                  │
          │    sec/op    │   sec/op     vs base                │    sec/op     vs base                │
And-12      354.88µ ± 1%   84.98µ ± 1%  -76.05% (p=0.000 n=10)   26.38µ ±  1%  -92.57% (p=0.000 n=10)
Or-12       357.74µ ± 2%   84.90µ ± 2%  -76.27% (p=0.000 n=10)   26.23µ ±  1%  -92.67% (p=0.000 n=10)
Xor-12      365.42µ ± 3%   85.10µ ± 1%  -76.71% (p=0.000 n=10)   26.38µ ±  1%  -92.78% (p=0.000 n=10)
AndNot-12   358.92µ ± 1%   85.65µ ± 4%  -76.14% (p=0.000 n=10)   26.47µ ±  0%  -92.63% (p=0.000 n=10)
Memset-12   297.70µ ± 1%   80.65µ ± 1%  -72.91% (p=0.000 n=10)   12.27µ ± 12%  -95.88% (p=0.000 n=10)
Popcnt-12   293.81µ ± 1%   45.94µ ± 1%  -84.36% (p=0.000 n=10)   11.42µ ±  1%  -96.11% (p=0.000 n=10)
geomean      336.7µ        76.14µ       -77.38%                  20.19µ        -94.00%

          │    naive     │                 purego                 │                   asm                    │
          │     B/s      │      B/s       vs base                 │      B/s        vs base                  │
And-12      2.624Gi ± 1%   10.959Gi ± 1%  +317.58% (p=0.000 n=10)   35.305Gi ±  1%  +1245.28% (p=0.000 n=10)
Or-12       2.603Gi ± 2%   10.970Gi ± 2%  +321.35% (p=0.000 n=10)   35.502Gi ±  1%  +1263.65% (p=0.000 n=10)
Xor-12      2.549Gi ± 3%   10.944Gi ± 1%  +329.42% (p=0.000 n=10)   35.309Gi ±  1%  +1285.41% (p=0.000 n=10)
AndNot-12   2.595Gi ± 1%   10.873Gi ± 4%  +319.04% (p=0.000 n=10)   35.183Gi ±  0%  +1255.91% (p=0.000 n=10)
Memset-12   3.128Gi ± 1%   11.548Gi ± 1%  +269.12% (p=0.000 n=10)   75.993Gi ± 13%  +2329.13% (p=0.000 n=10)
Popcnt-12   3.170Gi ± 1%   20.271Gi ± 1%  +539.48% (p=0.000 n=10)   81.560Gi ±  1%  +2472.98% (p=0.000 n=10)
geomean     2.766Gi         12.23Gi       +342.17%                   46.14Gi        +1568.02%
```
