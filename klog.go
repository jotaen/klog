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

// BinaryVersion is the version identifier that is embedded in the binary.
// It should be provided via a compile/build flag. The value must start with a
// `v`, followed by 2 dot-separated integers; e.g., `v2.12`. Optionally, there
// can be a dot- or hyphen-separated suffix (containing arbitrary characters),
// e.g., `v12.9-alpha` or `v8.1.0-alpha_2*`. If the build flag is absent, a value
// of `v?.?` is used.
var BinaryVersion string

// BinaryBuildHash is a computed hash of the underlying sources of the build.
// It should be provided via a compile/build flag. This is usually the git commit
// hash, e.g. `7ab819c`. The value must consist of at least 7 hexadecimal digits.
// Longer values are truncated; the value will also be normalised to lower-case. If
// the build flag is absent, a value of `???????` is used.
var BinaryBuildHash string

// BinaryBuildDate denotes the date that the build refers to. It should be provided
// via a compile/build flag. This should be the date of the latest commit, or
// alternatively the date at which the build occurred. It must be a valid calendar
// date, formatted `YYYY-MM-DD`; e.g. `2026-10-26`. If the build flag is absent, a
// value of `????-??-??` is used.
var BinaryBuildDate string

func main() {

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

	code, err := klog.Run(
		klogFolder,
		app.NewMeta(specification, license, BinaryVersion, BinaryBuildHash, BinaryBuildDate),
		config,
		os.Args[1:],
	)
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
