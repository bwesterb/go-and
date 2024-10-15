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
          │     naive     │                purego                │                 asm                 │
          │    sec/op     │    sec/op     vs base                │   sec/op     vs base                │
And-32       8162.0n ± 5%   2034.5n ± 1%  -75.07% (p=0.000 n=10)   631.2n ± 2%  -92.27% (p=0.000 n=10)
Or-32        9751.5n ± 8%   2104.5n ± 3%  -78.42% (p=0.000 n=10)   626.4n ± 1%  -93.58% (p=0.000 n=10)
Xor-32       8112.5n ± 3%   2029.0n ± 0%  -74.99% (p=0.000 n=10)   631.6n ± 1%  -92.22% (p=0.000 n=10)
AndNot-32   10685.5n ± 4%   2292.0n ± 2%  -78.55% (p=0.000 n=10)   635.2n ± 2%  -94.06% (p=0.000 n=10)
Memset-32    167.96µ ± 0%    57.54µ ± 1%  -65.74% (p=0.000 n=10)   15.83µ ± 1%  -90.57% (p=0.000 n=10)
Popcnt-32    132.15µ ± 1%    71.63µ ± 1%  -45.80% (p=0.000 n=10)   33.86µ ± 6%  -74.38% (p=0.000 n=10)
geomean       23.13µ         6.592µ       -71.50%                  2.097µ       -90.93%

          │    naive     │                 purego                 │                   asm                   │
          │     B/s      │      B/s       vs base                 │      B/s       vs base                  │
And-32      3.651Gi ± 5%   14.649Gi ± 1%  +301.20% (p=0.000 n=10)   47.212Gi ± 2%  +1193.01% (p=0.000 n=10)
Or-32       3.057Gi ± 8%   14.163Gi ± 3%  +363.37% (p=0.000 n=10)   47.580Gi ± 1%  +1456.63% (p=0.000 n=10)
Xor-32      3.674Gi ± 3%   14.690Gi ± 0%  +299.88% (p=0.000 n=10)   47.190Gi ± 1%  +1184.58% (p=0.000 n=10)
AndNot-32   2.789Gi ± 4%   13.003Gi ± 2%  +366.21% (p=0.000 n=10)   46.916Gi ± 2%  +1582.18% (p=0.000 n=10)
Memset-32   5.545Gi ± 0%   16.187Gi ± 1%  +191.91% (p=0.000 n=10)   58.816Gi ± 1%   +960.69% (p=0.000 n=10)
Popcnt-32   7.048Gi ± 1%   13.002Gi ± 1%   +84.48% (p=0.000 n=10)   27.506Gi ± 6%   +290.28% (p=0.000 n=10)
geomean     4.058Gi         14.24Gi       +250.89%                   44.76Gi       +1002.97%
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
