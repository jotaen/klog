package cli

import (
	"github.com/jotaen/klog/klog/app"
)

const DESCRIPTION = "klog: command line app for time tracking with plain-text files.\n" +
	"Run with --help to learn usage.\n" +
	"Documentation online at " + KLOG_WEBSITE_URL

type Info struct {
	Version bool `short:"v" name:"version" help:"Alias for 'klog version'"`
	Spec    bool `name:"spec" help:"Print file format specification"`
	License bool `name:"license" help:"Print license"`
}

func (opt *Info) Help() string {
	return DESCRIPTION
}

func (opt *Info) Run(ctx app.Context) app.Error {
	if opt.Version {
		versionCmd := Version{}
		return versionCmd.Run(ctx)
	} else if opt.Spec {
		ctx.Print(ctx.Meta().Specification + "\n")
		return nil
	} else if opt.License {
		ctx.Print(ctx.Meta().License + "\n")
		return nil
	}
	ctx.Print(DESCRIPTION + "\n")
	return nil
}
