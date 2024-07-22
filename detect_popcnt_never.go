//go:build !amd64 || purego

package and

func hasPopcnt() bool {
	return false
}
