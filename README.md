go-and
======

Fast bitwise and, or and andn for `[]byte` slices.

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

```
goos: linux
goarch: amd64
pkg: github.com/bwesterb/go-and
cpu: 13th Gen Intel(R) Core(TM) i9-13900
          │    naive     │                purego                │                 asm                 │
          │    sec/op    │    sec/op     vs base                │   sec/op     vs base                │
And-32      273.05µ ± 5%    64.48µ ± 2%  -76.39% (p=0.000 n=10)   21.88µ ± 1%  -91.99% (p=0.000 n=10)
Or-32       274.70µ ± 6%    64.36µ ± 1%  -76.57% (p=0.000 n=10)   21.81µ ± 1%  -92.06% (p=0.000 n=10)
AndNot-32   310.78µ ± 2%    71.01µ ± 2%  -77.15% (p=0.000 n=10)   21.83µ ± 1%  -92.98% (p=0.000 n=10)
Memset-32   167.77µ ± 0%   167.55µ ± 0%   -0.13% (p=0.002 n=10)   15.88µ ± 1%  -90.53% (p=0.000 n=10)
Popcnt-32   126.84µ ± 0%    71.42µ ± 1%  -43.69% (p=0.000 n=10)   32.48µ ± 1%  -74.40% (p=0.000 n=10)
geomean      218.3µ         81.18µ       -62.82%                  22.18µ       -89.84%

          │    naive     │                 purego                 │                   asm                   │
          │     B/s      │      B/s       vs base                 │      B/s       vs base                  │
And-32      3.411Gi ± 5%   14.444Gi ± 2%  +323.45% (p=0.000 n=10)   42.560Gi ± 1%  +1147.72% (p=0.000 n=10)
Or-32       3.391Gi ± 7%   14.470Gi ± 1%  +326.78% (p=0.000 n=10)   42.708Gi ± 1%  +1159.61% (p=0.000 n=10)
AndNot-32   2.997Gi ± 2%   13.116Gi ± 2%  +337.68% (p=0.000 n=10)   42.665Gi ± 1%  +1323.72% (p=0.000 n=10)
Memset-32   5.551Gi ± 0%    5.559Gi ± 0%    +0.13% (p=0.002 n=10)   58.642Gi ± 1%   +956.36% (p=0.000 n=10)
Popcnt-32   7.342Gi ± 0%   13.040Gi ± 1%   +77.60% (p=0.000 n=10)   28.677Gi ± 1%   +290.57% (p=0.000 n=10)
geomean     4.266Gi         11.47Gi       +168.93%                   41.98Gi        +884.14%
```
