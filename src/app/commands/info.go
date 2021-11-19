package commands

import (
	"github.com/jotaen/klog/src/app"
)

const DESCRIPTION = "klog: command line app for time tracking with plain-text files.\n" +
	"Run with --help to learn usage.\n" +
	"Documentation online at https://klog.jotaen.net"

type Info struct {
	Version bool `short:"v" name:"version" help:"Alias for 'klog version'"`
}

func (opt *Info) Run(ctx app.Context) error {
	if opt.Version {
		versionCmd := Version{}
		return versionCmd.Run(ctx)
	}
	ctx.Print(DESCRIPTION + "\n")
	return nil
}
