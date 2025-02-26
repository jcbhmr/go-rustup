package exemanip

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"golang.org/x/sys/windows"
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
	exeOldUTF16Ptr, err := windows.UTF16PtrFromString(exeOld)
	if err != nil {
		panic(err)
	}
	err = windows.SetFileAttributes(exeOldUTF16Ptr, windows.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		log.Println(err)
	}

	cmd := &exec.Cmd{
		Path:   exe,
		Args:   os.Args,
		Env:    os.Environ(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	signal.Ignore(os.Interrupt)
	signal.Ignore(os.Kill)
	err = cmd.Run()
	if err != nil {
		var exitError *exec.ExitError
		if !errors.As(err, &exitError) {
			log.Fatal(err)
		}
	}
	cmd2 := exec.Command("cmd.exe")
	cmd2.SysProcAttr = &windows.SysProcAttr{
		CmdLine:    "cmd.exe /C choice /C Y /N /D Y /T 3 > NUL & del " + exeOld,
		HideWindow: true,
	}
	err = cmd2.Start()
	if err != nil {
		log.Println(err)
	}
	os.Exit(cmd.ProcessState.ExitCode())
}
