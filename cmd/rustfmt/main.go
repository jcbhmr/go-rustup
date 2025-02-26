package main

import (
	"github.com/jcbhmr/go-rustup/internal/exemanip"
	"github.com/jcbhmr/go-rustup/internal/rustupinit"
)

func main() {
	exemanip.MoltMain(rustupinit.ExecutableBytes())
}
