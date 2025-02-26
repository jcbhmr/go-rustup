package xruntime

import (
	"os"
	"sync"
)

var ABI = sync.OnceValue(func() string {
	// Unsure if this is correct but it seems to work so far.
	_, ok := os.LookupEnv("MSYSTEM")
	if ok {
		return "gnu"
	} else {
		return "msvc"
	}
})
