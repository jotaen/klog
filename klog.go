package main

import (
	_ "embed"
	"fmt"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/main"
	"os"
	"os/user"
	"runtime"
)

//go:embed Specification.md
var specification string

//go:embed LICENSE.txt
var license string

//go:embed CHANGELOG.md
var changelog string

var BinaryVersion string   // Set via build flag
var BinaryBuildHash string // Set via build flag

func main() {
	if len(BinaryBuildHash) > 7 {
		BinaryBuildHash = BinaryBuildHash[:7]
	}
	prefs := app.NewDefaultPreferences()
	if os.Getenv("KLOG_DEBUG") != "" {
		prefs.IsDebug = true
	}
	if os.Getenv("NO_COLOR") != "" {
		prefs.NoColour = true
	}
	if os.Getenv("KLOG_BETA_PARALLEL") != "" {
		prefs.CpuKernels = runtime.NumCPU()
	}
	prefs.Editor = os.Getenv("KLOG_EDITOR")
	if prefs.Editor == "" {
		prefs.Editor = os.Getenv("EDITOR")
	}
	homeDir, err := user.Current()
	if err != nil {
		fmt.Println("Failed to initialise application. Error:")
		fmt.Println(err)
		os.Exit(1)
	}
	exitCode, runErr := klog.Run(homeDir.HomeDir, app.Meta{
		Specification: specification,
		License:       license,
		Changelog:     changelog,
		Version:       BinaryVersion,
		SrcHash:       BinaryBuildHash,
	}, prefs, os.Args[1:])
	if runErr != nil {
		fmt.Println(runErr)
	}
	os.Exit(exitCode)
}
