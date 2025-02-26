//go:build unix

package exemanip

import (
	"log"
	"os"

	"golang.org/x/sys/unix"
)

// Replaces the current on-disk executable with the given bytes and relaunches the process using the new executable.
//
// Returns never. Exits on error.
func MoltMain(bytes []byte) {
	exe, err := Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeInfo, err := os.Stat(exe)
	if err != nil {
		log.Fatal(err)
	}

	exeNew := exe + ".new"
	err = os.WriteFile(exeNew, bytes, exeInfo.Mode())
	if err != nil {
		log.Fatal(err)
	}

	exeOld, err := renamePreserve(exeNew, exe, "")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(exeOld)
	if err != nil {
		log.Println(err)
	}

	log.Fatal(unix.Exec(exe, os.Args, os.Environ()))
}
