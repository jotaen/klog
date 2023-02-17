package cli

import (
	"github.com/jotaen/klog/klog/app"
	"path/filepath"
	"strings"
)

const DESCRIPTION = "klog: command line app for time tracking with plain-text files.\n" +
	"Run with --help to learn usage.\n" +
	"Documentation online at " + KLOG_WEBSITE_URL

type Info struct {
	Version      bool `short:"v" name:"version" help:"Alias for 'klog version'"`
	ConfigFolder bool `name:"config-folder" help:"Prints path of klog config folder"`
	Spec         bool `name:"spec" help:"Prints file format specification"`
	License      bool `name:"license" help:"Prints license"`
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
	} else if opt.ConfigFolder {
		ctx.Print(ctx.KlogConfigFolder().Path() + "\n")
		lookups := make([]string, len(app.KLOG_CONFIG_FOLDER))
		for i, f := range app.KLOG_CONFIG_FOLDER {
			lookups[i] = filepath.Join(f.EnvVarSymbol(), f.Location)
		}
		ctx.Print("(Lookup order: " + strings.Join(lookups, ", ") + ")\n")
		return nil
	}
	ctx.Print(DESCRIPTION + "\n")
	return nil
}
