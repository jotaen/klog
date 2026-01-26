//go:build darwin

package app

import (
	"github.com/jotaen/klog/lib/shellcmd"
)

var POTENTIAL_EDITORS = []shellcmd.Command{
	shellcmd.New("vim", nil),
	shellcmd.New("vi", nil),
	shellcmd.New("nano", nil),
	shellcmd.New("pico", nil),
	shellcmd.New("open", []string{"-a", "TextEdit"}),
}

var POTENTIAL_FILE_EXLORERS = []shellcmd.Command{
	shellcmd.New("open", nil),
}

var KLOG_CONFIG_FOLDER = []KlogFolder{
	{"KLOG_CONFIG_HOME", ""},
	{"XDG_CONFIG_HOME", "klog"},
	{"HOME", ".klog"},
}

func (kf KlogFolder) EnvVarSymbol() string {
	return "$" + kf.BasePathEnvVar
}
