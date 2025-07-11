package cli

import (
	"github.com/jotaen/klog/klog/app"
)

type Info struct {
	Spec    bool `name:"spec" help:"Print the .klg file format specification."`
	License bool `name:"license" help:"Print license / copyright information."`
	About   bool `name:"about" help:"Print meta information about klog."`
}

func (opt *Info) Run(ctx app.Context) app.Error {
	//ctx.Print(ctx.KlogConfigFolder().Path() + "\n")
	if opt.Spec {
		ctx.Print(ctx.Meta().Specification + "\n")
	} else if opt.License {
		ctx.Print(ctx.Meta().License + "\n")
	} else if opt.About {
		ctx.Print("klog is a " + INTRO_SUMMARY + "\n")
	} else {
		ctx.Print("Use --spec or --license\n")
	}
	return nil
}
