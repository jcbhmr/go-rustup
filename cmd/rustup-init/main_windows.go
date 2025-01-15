package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"

	"golang.org/x/sys/windows"
)

//go:generate go run ./gen.go

// 1. Get properties of rustup-init (self)
// 2. Write rustup-init (new) to rustup-init.NEW with the same properties.
// 3. Remove any existing rustup-init.OLD file.
// 4. Move rustup-init (self) to rustup-init.OLD.
// 5. Move rustup-init.NEW to rustup-init (new).
// 6. Mark rustup-init.OLD as hidden. You cannot delete your own executable on Windows.
// 7. Mark rustup-init.OLD to be deleted when rebooted.
// 8. Run rustup-init (new) with the same arguments.
// 9. Launch an orphaned shell to delete rustup-init.OLD after this process exits.
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
		// log.Printf("warning: %v", err)
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
	exeOldUTF16Ptr, err := windows.UTF16PtrFromString(exe + ".OLD")
	if err != nil {
		panic(err)
	}
	err = windows.SetFileAttributes(exeOldUTF16Ptr, windows.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		// log.Printf("warning: %v", err)
	}
	cmd := &exec.Cmd{
		Path:   exe,
		Args:   os.Args,
		Env:    os.Environ(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err = cmd.Run()
	if exitError := (*exec.ExitError)(nil); errors.As(err, &exitError) {
		// Continue
	} else if err != nil {
		log.Fatal(err)
	}
	cmd2 := exec.Command("cmd.exe", "/C choice /C Y /N /D Y /T 3 > NUL & del \""+exe+".OLD\"")
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	err = cmd2.Start()
	if err != nil {
		// log.Printf("warning: %v", err)
	}
	os.Exit(cmd.ProcessState.ExitCode())
}
