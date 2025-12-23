//go:build windows

package app

import (
	"github.com/jotaen/klog/lib/shellcmd"

	"os"
	"syscall"
	"unsafe"
)

var POTENTIAL_EDITORS = []shellcmd.Command{
	shellcmd.New("notepad", nil),
}

var POTENTIAL_FILE_EXLORERS = []shellcmd.Command{
	shellcmd.New("cmd.exe", []string{"/C", "start"}),
}

var KLOG_CONFIG_FOLDER = []KlogFolder{
	{"KLOG_CONFIG_HOME", ""},
	{"AppData", "klog"},
}

func (kf KlogFolder) EnvVarSymbol() string {
	return "%" + kf.BasePathEnvVar + "%"
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
