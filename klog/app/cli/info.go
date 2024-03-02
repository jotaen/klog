package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
	"path/filepath"
	"strings"
)

type Info struct {
	Spec         InfoSpec         `cmd:"" name:"spec" help:"Prints file format specification"`
	License      InfoLicense      `cmd:"" name:"license" help:"Prints license / copyright information"`
	ConfigFolder InfoConfigFolder `cmd:"" name:"config-folder" help:"Prints path of klog config folder"`
}

func (opt *Info) Help() string {
	return ""
}

type InfoConfigFolder struct {
	util.QuietArgs
}

func (opt *InfoConfigFolder) Run(ctx app.Context) app.Error {
	ctx.Print(ctx.KlogConfigFolder().Path() + "\n")
	if !opt.Quiet {
		lookups := make([]string, len(app.KLOG_CONFIG_FOLDER))
		for i, f := range app.KLOG_CONFIG_FOLDER {
			lookups[i] = filepath.Join(f.EnvVarSymbol(), f.Location)
		}
		ctx.Print("(Lookup order: " + strings.Join(lookups, ", ") + ")\n")
	}
	return nil
}

type InfoSpec struct{}

func (opt *InfoSpec) Run(ctx app.Context) app.Error {
	ctx.Print(ctx.Meta().Specification + "\n")
	return nil
}

type InfoLicense struct{}

func (opt *InfoLicense) Run(ctx app.Context) app.Error {
	ctx.Print(ctx.Meta().License + "\n")
	return nil
}
