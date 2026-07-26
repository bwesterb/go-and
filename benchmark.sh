#!/bin/sh

set -ex

FUNCTIONS="(And|Or|Xor|AndNot|Any|AnyMasked|Popcnt|PopcntMasked)"

go test -run=^# -count=10 -bench="^Benchmark${FUNCTIONS}$" | tee asm
go test -run=^# -count=10 -bench="^Benchmark${FUNCTIONS}$" -tags purego | tee purego
go test -run=^# -count=10 -bench="^Benchmark${FUNCTIONS}Naive$" | sed --unbuffered 's/Naive//g' | tee naive
go run golang.org/x/perf/cmd/benchstat@latest naive purego asm
