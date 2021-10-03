//go:build !windows
// +build !windows

package app

var POTENTIAL_EDITORS = []string{"vim", "vi", "nano", "pico"}

func flagAsHidden(path string) {
	// Nothing to do on UNIX due to the dotfile convention
}
