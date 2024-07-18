//go:build !amd64

package and

func hasPopcnt() bool {
	return false
}
