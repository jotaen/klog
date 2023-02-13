package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
)

type Config struct {
	ConfigFilePath bool `name:"file-path" help:"Prints the path to your config file"`
}

func (opt *Config) Help() string {
	return `You are able to configure some of klog’s behaviour via a YAML file in your ` + "`" + app.KLOG_FOLDER_NAME + "`" + ` folder. (Run ` + "`" + `klog config --file-path` + "`" + ` to print the exact location.)

If you run ` + "`" + `klog config` + "`" + `, you can learn about the supported YAML properties in the file, and you also see what values are in effect at the moment.

Note: the output of the command does not print the actual file. You may, however, use the output as template for setting up the file, as its YAML-formatted.`
}

func (opt *Config) Run(ctx app.Context) app.Error {
	if opt.ConfigFilePath {
		ctx.Print(app.Join(ctx.KlogFolder(), app.CONFIG_FILE_NAME).Path() + "\n")
		return nil
	}
	for i, e := range app.CONFIG_FILE_ENTRIES {
		ctx.Print(lib.Subdued.Format(lib.Reflower.Reflow(e.Description+"\n"+e.Instructions, "# ")))
		ctx.Print("\n")
		ctx.Print(lib.Red.Format(e.Name) + `: ` + terminalformat.Style{Color: "227"}.Format(e.Value(ctx.Config())))
		if i < len(app.CONFIG_FILE_ENTRIES)-1 {
			ctx.Print("\n\n")
		}
	}
	ctx.Print("\n")
	return nil
}
