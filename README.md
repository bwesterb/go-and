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
And-32            6942.5n ± 2%   2825.5n ± 0%  -59.30% (p=0.000 n=10)   360.8n ± 1%  -94.80% (p=0.000 n=10)
Or-32             6858.0n ± 1%   2839.5n ± 0%  -58.60% (p=0.000 n=10)   365.7n ± 3%  -94.67% (p=0.000 n=10)
Xor-32            6996.5n ± 2%   2829.0n ± 0%  -59.57% (p=0.000 n=10)   361.9n ± 2%  -94.83% (p=0.000 n=10)
AndNot-32         7597.0n ± 1%   2836.5n ± 0%  -62.66% (p=0.000 n=10)   360.4n ± 1%  -95.26% (p=0.000 n=10)
AnyMasked-32      5722.0n ± 1%   1337.5n ± 0%  -76.63% (p=0.000 n=10)   506.6n ± 3%  -91.15% (p=0.000 n=10)
Any-32            5678.5n ± 0%   1072.5n ± 1%  -81.11% (p=0.000 n=10)   261.1n ± 0%  -95.40% (p=0.000 n=10)
Popcnt-32         177.46µ ± 1%    66.45µ ± 1%  -62.55% (p=0.000 n=10)   22.49µ ± 1%  -87.33% (p=0.000 n=10)
PopcntMasked-32   315.53µ ± 0%   115.54µ ± 3%  -63.38% (p=0.000 n=10)   28.07µ ± 8%  -91.10% (p=0.000 n=10)
geomean            16.14µ         5.387µ       -66.62%                  1.046µ       -93.52%

                │    naive     │                 purego                 │                   asm                    │
                │     B/s      │      B/s       vs base                 │      B/s        vs base                  │
And-32            4.293Gi ± 1%   10.547Gi ± 0%  +145.69% (p=0.000 n=10)    82.618Gi ± 1%  +1824.54% (p=0.000 n=10)
Or-32             4.345Gi ± 1%   10.495Gi ± 0%  +141.52% (p=0.000 n=10)    81.486Gi ± 3%  +1775.23% (p=0.000 n=10)
Xor-32            4.259Gi ± 1%   10.534Gi ± 0%  +147.31% (p=0.000 n=10)    82.352Gi ± 2%  +1833.41% (p=0.000 n=10)
AndNot-32         3.923Gi ± 1%   10.505Gi ± 0%  +167.80% (p=0.000 n=10)    82.691Gi ± 1%  +2007.97% (p=0.000 n=10)
AnyMasked-32      5.208Gi ± 1%   22.278Gi ± 0%  +327.74% (p=0.000 n=10)    58.832Gi ± 3%  +1029.57% (p=0.000 n=10)
Any-32            5.248Gi ± 0%   27.792Gi ± 1%  +429.53% (p=0.000 n=10)   114.156Gi ± 0%  +2075.03% (p=0.000 n=10)
Popcnt-32         5.248Gi ± 1%   14.015Gi ± 1%  +167.05% (p=0.000 n=10)    41.405Gi ± 1%   +688.95% (p=0.000 n=10)
PopcntMasked-32   2.952Gi ± 0%    8.061Gi ± 3%  +173.09% (p=0.000 n=10)    33.178Gi ± 7%  +1024.07% (p=0.000 n=10)
geomean           4.366Gi         13.08Gi       +199.60%                    67.34Gi       +1442.50%
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
