//go:build darwin

package app

import (
	"github.com/jotaen/klog/klog/app/cli/lib/command"
)

var POTENTIAL_EDITORS = []command.Command{
	command.New("vim", nil),
	command.New("vi", nil),
	command.New("nano", nil),
	command.New("pico", nil),
	command.New("open", []string{"-a", "TextEdit"}),
}

var POTENTIAL_FILE_EXLORERS = []command.Command{
	command.New("open", nil),
}

var KLOG_CONFIG_FOLDER = []KlogFolder{
	{"KLOG_CONFIG_HOME", ""},
	{"XDG_CONFIG_HOME", "klog"},
	{"HOME", ".klog"},
}

func (kf KlogFolder) EnvVarSymbol() string {
	return "$" + kf.BasePathEnvVar
}
