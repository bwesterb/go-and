//go:build amd64.v2 && !purego

package and

func hasPopcnt() bool {
	return true
}
