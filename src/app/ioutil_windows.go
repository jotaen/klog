//go:build windows
// +build windows

package app

import (
	"os"
	"syscall"
)

var POTENTIAL_EDITORS = []string{"notepad"}

func flagAsHidden(path string) {
	winFileName, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return
	}
	_ = syscall.SetFileAttributes(winFileName, syscall.FILE_ATTRIBUTE_HIDDEN)
}

func createSymlinkForBookmark(targetPath string, linkPath string) Error {
	err := os.Symlink(targetPath, linkPath)
	if err != nil {
		return NewError(
			"Failed to create bookmark",
			"On Windows the `bookmark set` subcommand requires admin privileges or "+
				"developer mode. (klog needs to create a symlink and Windows doesn’t "+
				"allow this otherwise.)",
			err,
		)
	}
	return nil
}
