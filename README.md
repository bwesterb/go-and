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
And-32      8160.5n ± 2%   1576.5n ± 1%  -80.68% (p=0.000 n=10)   524.9n ± 2%  -93.57% (p=0.000 n=10)
Or-32       8132.5n ± 3%   1798.0n ± 1%  -77.89% (p=0.000 n=10)   520.1n ± 1%  -93.60% (p=0.000 n=10)
Xor-32      8320.0n ± 2%   1547.5n ± 0%  -81.40% (p=0.000 n=10)   516.0n ± 1%  -93.80% (p=0.000 n=10)
AndNot-32   9922.5n ± 1%   2062.0n ± 1%  -79.22% (p=0.000 n=10)   516.6n ± 2%  -94.79% (p=0.000 n=10)
Memset-32   257.22µ ± 8%    35.86µ ± 3%  -86.06% (p=0.000 n=10)   15.81µ ± 1%  -93.85% (p=0.000 n=10)
Not-32      272.72µ ± 2%    44.50µ ± 1%  -83.68% (p=0.000 n=10)   15.26µ ± 1%  -94.41% (p=0.000 n=10)
Popcnt-32   129.13µ ± 1%    51.45µ ± 2%  -60.16% (p=0.000 n=10)   35.93µ ± 3%  -72.18% (p=0.000 n=10)
geomean      33.73µ         6.897µ       -79.55%                  2.512µ       -92.55%

          │    naive     │                 purego                 │                   asm                   │
          │     B/s      │      B/s       vs base                 │      B/s       vs base                  │
And-32      3.652Gi ± 2%   18.903Gi ± 1%  +417.63% (p=0.000 n=10)   56.768Gi ± 2%  +1454.46% (p=0.000 n=10)
Or-32       3.665Gi ± 3%   16.577Gi ± 1%  +352.36% (p=0.000 n=10)   57.301Gi ± 1%  +1463.65% (p=0.000 n=10)
Xor-32      3.582Gi ± 2%   19.257Gi ± 0%  +437.58% (p=0.000 n=10)   57.755Gi ± 1%  +1512.31% (p=0.000 n=10)
AndNot-32   3.004Gi ± 1%   14.452Gi ± 1%  +381.16% (p=0.000 n=10)   57.693Gi ± 2%  +1820.84% (p=0.000 n=10)
Memset-32   3.621Gi ± 9%   25.973Gi ± 3%  +617.24% (p=0.000 n=10)   58.911Gi ± 1%  +1526.84% (p=0.000 n=10)
Not-32      3.415Gi ± 2%   20.929Gi ± 1%  +512.84% (p=0.000 n=10)   61.036Gi ± 1%  +1687.27% (p=0.000 n=10)
Popcnt-32   7.212Gi ± 1%   18.101Gi ± 2%  +150.97% (p=0.000 n=10)   25.923Gi ± 3%   +259.44% (p=0.000 n=10)
geomean     3.863Gi         18.89Gi       +388.95%                   51.87Gi       +1242.68%
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
