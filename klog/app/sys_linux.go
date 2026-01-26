//go:build linux

package app

import (
	"github.com/jotaen/klog/lib/shellcmd"
)

var POTENTIAL_EDITORS = []shellcmd.Command{
	shellcmd.New("vim", nil),
	shellcmd.New("vi", nil),
	shellcmd.New("nano", nil),
	shellcmd.New("pico", nil),
}

var POTENTIAL_FILE_EXLORERS = []shellcmd.Command{
	shellcmd.New("xdg-open", nil),
}

var KLOG_CONFIG_FOLDER = []KlogFolder{
	{"KLOG_CONFIG_HOME", ""},
	{"XDG_CONFIG_HOME", "klog"},
	{"HOME", ".config/klog"},
}

func (kf KlogFolder) EnvVarSymbol() string {
	return "$" + kf.BasePathEnvVar
}
