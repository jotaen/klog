package main

import (
	_ "embed"
	"fmt"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/main"
	"os"
	"os/user"
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
	isDebug := false
	if os.Getenv("KLOG_DEBUG") != "" {
		isDebug = true
	}
	homeDir, err := user.Current()
	if err != nil {
		fmt.Println("Failed to initialise application. Error:")
		fmt.Println(err)
		os.Exit(1)
	}
	runErr, exitCode := klog.Run(homeDir.HomeDir, app.Meta{
		Specification: specification,
		License:       license,
		Version:       BinaryVersion,
		BuildHash:     BinaryBuildHash,
	}, isDebug, os.Args[1:])
	if runErr != nil {
		fmt.Println(runErr)
	}
	os.Exit(exitCode)
}
