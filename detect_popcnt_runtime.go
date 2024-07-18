//go:build amd64 && !amd64.v2

package and

import "golang.org/x/sys/cpu"

func hasPopcnt() bool {
	return cpu.X86.HasPOPCNT
}
