package rustupinit

import (
	_ "embed"

	"github.com/jcbhmr/go-rustup/internal/ezgzip"
	"github.com/jcbhmr/go-rustup/internal/xruntime"
)

//go:embed rustup-init.windows-amd64-gnu.exe.gz
var gnu []byte

//go:embed rustup-init.windows-amd64-msvc.exe.gz
var msvc []byte

func gzippedExecutableBytes() []byte {
	if xruntime.ABI() == "gnu" {
		return gnu
	} else {
		return msvc
	}
}

func ExecutableBytes() []byte {
	return ezgzip.MustDecompressBytes(gzippedExecutableBytes())
}
