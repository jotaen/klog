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

var BinaryVersion string   // Set via build flag
var BinaryBuildHash string // Set via build flag

func main() {
	if len(BinaryBuildHash) > 7 {
		BinaryBuildHash = BinaryBuildHash[:7]
	}

	homeDir, err := user.Current()
	if err != nil {
		fmt.Println("Failed to initialise application. Error:")
		fmt.Println(err)
		os.Exit(1)
	}

	config := app.NewConfig(
		app.ConfigFromStaticValues{NumCpus: runtime.NumCPU()},
		app.ConfigFromEnvVars{GetVar: os.Getenv},
	)

	exitCode, runErr := klog.Run(homeDir.HomeDir, app.Meta{
		Specification: specification,
		License:       license,
		Version:       BinaryVersion,
		SrcHash:       BinaryBuildHash,
	}, config, os.Args[1:])
	if runErr != nil {
		fmt.Println(runErr)
	}
	os.Exit(exitCode)
}
