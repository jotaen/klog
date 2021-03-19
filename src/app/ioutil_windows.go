// +build windows

package app

import (
    "syscall"
)

func flagAsHidden(path string) {
	winFileName, err := syscall.UTF16PtrFromString(path)
    if err != nil {
        return
    }
    _ = syscall.SetFileAttributes(winFileName, syscall.FILE_ATTRIBUTE_HIDDEN)
}
