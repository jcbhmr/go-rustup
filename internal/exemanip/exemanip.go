package exemanip

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jcbhmr/go-rustup/internal/robustio"
)

// os.Executable() + filepath.EvalSymlinks()
func Executable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	res, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return "", err
	}
	return res, nil
}

// (Fallible os.Remove()) + robustio.Rename()
func renameOverwrite(oldpath, newpath string) error {
	err := os.Remove(newpath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("rename overwrite %q to %q: %w", oldpath, newpath, err)
	}
	err2 := robustio.Rename(oldpath, newpath)
	if err2 != nil {
		if err == nil {
			err = fmt.Errorf("already removed %q", newpath)
			return fmt.Errorf("rename overwrite %q to %q: %w; %w", oldpath, newpath, err, err2)
		}
		return fmt.Errorf("rename overwrite %q to %q: %w", oldpath, newpath, err2)
	}
	return nil
}

// Renames newpath to backuppath (overwriting if it exists), then renames oldpath to newpath. Does not remove backuppath when done.
//
// backuppath defaults to newpath + ".bak" if empty.
//
// Tries to revert everything on error.
func renamePreserve(oldpath, newpath, backuppath string) (string, error) {
	if backuppath == "" {
		backuppath = newpath + ".bak"
	}
	err := renameOverwrite(newpath, backuppath)
	if err != nil {
		return backuppath, fmt.Errorf("rename preserve %q to %q to %q: %w", oldpath, newpath, backuppath, err)
	}
	err = robustio.Rename(oldpath, newpath)
	if err != nil {
		err2 := robustio.Rename(backuppath, newpath)
		if err2 != nil {
			return backuppath, fmt.Errorf("rename preserve %q to %q to %q: %w; %w", oldpath, newpath, backuppath, err, err2)
		}
		return "", fmt.Errorf("rename preserve %q to %q to %q: %w", oldpath, newpath, backuppath, err)
	}
	return backuppath, nil
}
