package cli

import (
	"klog/app"
)

const DESCRIPTION = "klog: command line app for time tracking with plain-text files.\n" +
	"Run with --help to learn usage.\n" +
	"Documentation online at https://klog.jotaen.net"

type Info struct{}

func (opt *Info) Run(ctx app.Context, cli *Cli) error {
	if cli.VersionFlag {
		versionCmd := Version{}
		return versionCmd.Run(ctx)
	}
	ctx.Print(DESCRIPTION + "\n")
	return nil
}
