//go:build !windows
// +build !windows

package app

import "os"

var POTENTIAL_EDITORS = []string{"vim", "vi", "nano", "pico"}

func flagAsHidden(path string) {
	// Nothing to do on UNIX
}

func createSymlinkForBookmark(targetPath string, linkPath string) Error {
	err := os.Symlink(targetPath, linkPath)
	if err != nil {
		return NewError(
			"Failed to create bookmark",
			"Unable to create a symlink for the new bookmark",
			err,
		)
	}
	return nil
}
