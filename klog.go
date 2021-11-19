package main

import (
	_ "embed"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/main"
	"os"
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
	appErr := klog.Run(app.Meta{
		Specification: specification,
		License:       license,
		Version:       BinaryVersion,
		BuildHash:     BinaryBuildHash,
	})
	if appErr != nil {
		os.Exit(appErr.Code().ToInt())
	}
	os.Exit(0)
}
