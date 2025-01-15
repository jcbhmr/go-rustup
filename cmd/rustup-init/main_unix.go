//go:build unix

package main

import (
	"errors"
	"io/fs"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

//go:generate go run ./gen.go

// 1. Get properties of rustup-init (Go)
// 2. Write rustup-init (Rust) to rustup-init.NEW (Rust) with the same properties.
// 3. Remove any existing rustup-init.OLD file.
// 4. Move rustup-init (Go) to rustup-init.OLD (Go).
// 5. Move rustup-init.NEW (Rust) to rustup-init (Rust).
// 6. Remove rustup-init.OLD (Go).
// 7. exec*() the new rustup-init (Rust).
func main() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	oldFileInfo, err := os.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(exe+".NEW", rustupInit, oldFileInfo.Mode())
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(exe + ".OLD")
	if errors.Is(err, fs.ErrNotExist) {
		// Continue
	} else {
		log.Printf("warning: %v", err)
	}
	err = os.Rename(exe, exe+".OLD")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Rename(exe+".NEW", exe)
	if err != nil {
		err2 := os.Rename(exe+".OLD", exe)
		if err2 != nil {
			log.Fatalf("%v\n%v", err, err2)
		}
		log.Fatal(err)
	}
	err = os.Remove(exe + ".OLD")
	if err != nil {
		log.Printf("warning: %v", err)
	}
	log.Fatal(unix.Exec(exe, os.Args, os.Environ()))
}
