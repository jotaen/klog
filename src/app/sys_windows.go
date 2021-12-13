//go:build windows

package app

import (
	"syscall"
)

var POTENTIAL_EDITORS = []string{"notepad"}

var POTENTIAL_FILE_EXLORERS = []string{"start"}

func flagAsHidden(path string) {
	winFileName, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return
	}
	_ = syscall.SetFileAttributes(winFileName, syscall.FILE_ATTRIBUTE_HIDDEN)
}

func init() {
	enableAnsiEscapeSequences()
}

func enableAnsiEscapeSequences() {
	const enableVirtualTerminalProcessing = 0x0004

	var (
		kernel32           = syscall.NewLazyDLL("kernel32.dll")
		procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
		procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
	)

	var mode uint32
	procGetConsoleMode.Call(os.Stdout.Fd(), uintptr(unsafe.Pointer(&mode)))
	if (mode & enableVirtualTerminalProcessing) != enableVirtualTerminalProcessing {
		procSetConsoleMode.Call(os.Stdout.Fd(), uintptr(mode|enableVirtualTerminalProcessing))
	}
}
