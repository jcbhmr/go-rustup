//go:build !(linux || windows)

package xruntime

func ABI() string {
	return ""
}
