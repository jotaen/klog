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

	klogFolderPath := func() string {
		homeDir, hErr := user.Current()
		if hErr != nil {
			fail(hErr, false)
		}
		return homeDir.HomeDir + "/" + app.KLOG_FOLDER + "/"
	}()

	configFile := func() string {
		c, cErr := readConfigFile(klogFolderPath)
		if cErr != nil {
			fail(cErr, false)
		}
		return c
	}()

	config, cErr := app.NewConfig(
		app.FromStaticValues{NumCpus: runtime.NumCPU()},
		app.FromEnvVars{GetVar: os.Getenv},
		app.FromConfigFile{FileContents: configFile},
	)
	if cErr != nil {
		fail(cErr, false)
	}

	runErr := klog.Run(klogFolderPath, app.Meta{
		Specification: specification,
		License:       license,
		Version:       BinaryVersion,
		SrcHash:       BinaryBuildHash,
	}, config, os.Args[1:])
	if runErr != nil {
		fail(runErr, config.IsDebug.Value())
	}
}

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

func readConfigFile(klogFolderPath string) (string, error) {
	file, fErr := app.NewFile(klogFolderPath + app.CONFIG_FILE)
	if fErr != nil {
		return "", fErr
	}
	contents, rErr := app.ReadFile(file)
	if rErr != nil {
		if rErr.Code() == app.NO_SUCH_FILE {
			return "", nil
		}
		return "", rErr
	}
	return contents, nil
}
