//go:build linux

package app

import "github.com/jotaen/klog/klog/app/cli/lib/command"

var POTENTIAL_EDITORS = []command.Command{
	command.New("vim", nil),
	command.New("vi", nil),
	command.New("nano", nil),
	command.New("pico", nil),
}

var POTENTIAL_FILE_EXLORERS = []command.Command{
	command.New("xdg-open", nil),
}

func flagAsHidden(path string) {
	// Nothing to do on UNIX due to the dotfile convention
}
