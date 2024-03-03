package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
	"path/filepath"
	"strings"
)

type Info struct {
	Spec         InfoSpec         `cmd:"" name:"spec" help:"Print the .klg file format specification."`
	License      InfoLicense      `cmd:"" name:"license" help:"Print license / copyright information."`
	ConfigFolder InfoConfigFolder `cmd:"" name:"config-folder" help:"Print the path of the klog config folder."`
}

func (opt *Info) Help() string {
	return `
Run 'klog info config-folder' to see the location of your klog config folder.
The location of the config folder depends on your operating system and environment settings.
You can customise the folder’s location via environment variables – run the command to learn which ones klog relies on.

The config folder is used to store two files:
  - 'config.ini' (optional) – you can create this file manually to override some of klog’s default behaviour. Run 'klog config' to learn more.
  - 'bookmarks.json' (optional) – if you use bookmarks, then klog uses this file as database. You are not supposed to edit this file by hand! Instead, use the 'klog bookmarks' command to manage your bookmarks.

Run 'klog info spec' to read the formal specification of the klog file format.
If you want to review klog’s license and copyright information, run 'klog info license'.
`
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
