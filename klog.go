package main

import (
	_ "embed"
	"fmt"
	"os"
	"runtime"

	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/prettify"
	"github.com/jotaen/klog/klog/app/main"
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
		f, err := determineKlogConfigFolder()
		if err != nil {
			fail(prettify.PrettifyAppError(err, false), app.CONFIG_ERROR.ToInt())
		}
		return f
	}()

	configFile := func() string {
		c, err := readConfigFile(app.Join(klogFolder, app.CONFIG_FILE_NAME))
		if err != nil {
			fail(prettify.PrettifyAppError(err, false), app.CONFIG_ERROR.ToInt())
		}
		return c
	}()

	config := func() app.Config {
		c, err := app.NewConfig(runtime.NumCPU(), os.Getenv, configFile)
		if err != nil {
			fail(prettify.PrettifyAppError(err, false), app.CONFIG_ERROR.ToInt())
		}
		return c
	}()

	code, err := klog.Run(klogFolder, app.Meta{
		Specification: specification,
		License:       license,
		Version:       BinaryVersion,
		SrcHash:       BinaryBuildHash,
	}, config, os.Args[1:])
	if err != nil {
		fail(err, code)
	}
}

// fail terminates the process with an error.
func fail(err error, exitCode int) {
	fmt.Println(err)
	os.Exit(exitCode)
}

// readConfigFile reads the config file from disk, if present.
// If not present, it returns empty string.
func readConfigFile(location app.File) (string, app.Error) {
	contents, rErr := app.ReadFile(location)
	if rErr != nil {
		if rErr.Code() == app.NO_SUCH_FILE {
			return "", nil
		}
		return "", rErr
	}
	return contents, nil
}

// determineKlogConfigFolder returns the location where the klog config folder
// is (or should be) located.
func determineKlogConfigFolder() (app.File, app.Error) {
	for _, kf := range app.KLOG_CONFIG_FOLDER {
		basePath := os.Getenv(kf.BasePathEnvVar)
		if basePath != "" {
			return app.NewFile(basePath, kf.Location)
		}
	}
	return nil, app.NewError(
		"Cannot determine klog config folder",
		"Please set the `KLOG_CONFIG_HOME` environment variable, and make it point to a valid, empty folder.",
		nil,
	)
}
