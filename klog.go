package main

import (
	_ "embed"
	"fmt"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/main"
	"os"
	"os/user"
	"runtime"
)

//go:embed Specification.md
var specification string

//go:embed LICENSE.txt
var license string

var BinaryVersion string   // Set via build flag
var BinaryBuildHash string // Set via build flag

func main() {
	if len(BinaryBuildHash) > 7 {
		BinaryBuildHash = BinaryBuildHash[:7]
	}

	klogFolder := func() app.File {
		f, err := determineKlogFolderLocation()
		if err != nil {
			fail(err, false)
		}
		return app.Join(f, app.KLOG_FOLDER_NAME)
	}()

	configFile := func() string {
		c, err := readConfigFile(app.Join(klogFolder, app.CONFIG_FILE_NAME))
		if err != nil {
			fail(err, false)
		}
		return c
	}()

	config := func() app.Config {
		c, err := app.NewConfig(
			app.FromStaticValues{NumCpus: runtime.NumCPU()},
			app.FromEnvVars{GetVar: os.Getenv},
			app.FromConfigFile{FileContents: configFile},
		)
		if err != nil {
			fail(err, false)
		}
		return c
	}()

	err := klog.Run(klogFolder, app.Meta{
		Specification: specification,
		License:       license,
		Version:       BinaryVersion,
		SrcHash:       BinaryBuildHash,
	}, config, os.Args[1:])
	if err != nil {
		fail(err, config.IsDebug.Value())
	}
}

// fail terminates the process with an error.
func fail(e error, isDebug bool) {
	exitCode := -1
	if e != nil {
		fmt.Println(lib.PrettifyError(e, isDebug))
		if appErr, isAppError := e.(app.Error); isAppError {
			exitCode = appErr.Code().ToInt()
		}
	}
	os.Exit(exitCode)
}

// readConfigFile reads the config file from disk, if present.
// If not present, it returns empty string.
func readConfigFile(location app.File) (string, error) {
	contents, rErr := app.ReadFile(location)
	if rErr != nil {
		if rErr.Code() == app.NO_SUCH_FILE {
			return "", nil
		}
		return "", rErr
	}
	return contents, nil
}

// determineKlogFolderLocation returns the location where the `.klog` folder should be place.
// This is determined by following this lookup precedence:
// - $KLOG_FOLDER_LOCATION, if set
// - $XDG_CONFIG_HOME, if set
// - The default home folder, e.g. `~`
func determineKlogFolderLocation() (app.File, error) {
	location := os.Getenv("KLOG_FOLDER_LOCATION")
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		location = os.Getenv("XDG_CONFIG_HOME")
	} else if location == "" {
		homeDir, hErr := user.Current()
		if hErr != nil {
			return nil, hErr
		}
		location = homeDir.HomeDir
	}
	return app.NewFile(location)
}
