package xruntime

import (
	"bytes"
	"os"
	"os/exec"
	"sync"
)

var ABI = sync.OnceValue(func() string {
	ldd, err := os.ReadFile("/usr/bin/ldd")
	if err == nil {
		if bytes.Contains(ldd, []byte("musl")) {
			return "musl"
		}
		if bytes.Contains(ldd, []byte("GNU C Library")) {
			return "gnu"
		}
	}
	cmd := exec.Command("getconf", "GNU_LIBC_VERSION")
	if cmd.Err == nil {
		out, err := cmd.Output()
		if err == nil {
			if bytes.Contains(bytes.TrimSpace(out), []byte("glibc")) {
				return "gnu"
			}
		}
	}
	cmd = exec.Command("ldd", "--version")
	if cmd.Err == nil {
		out, err := cmd.Output()
		if err == nil {
			if bytes.Contains(bytes.TrimSpace(out), []byte("musl")) {
				return "musl"
			}
		}
	}
	return ""
})
