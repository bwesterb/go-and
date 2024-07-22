#!/bin/sh

set -ex

# ^.{2,6}$ is a hack to skip the .+Generic benchmarks
go test -run=^# -count=10 -bench="^Benchmark.{2,6}$" | tee asm
go test -run=^# -count=10 -bench="^Benchmark.{2,6}$" -tags purego | tee purego
go test -run=^# -count=10 -bench="^Benchmark.{2,6}Naive$" | sed --unbuffered 's/Naive//g' | tee naive
go run golang.org/x/perf/cmd/benchstat@latest naive purego asm
