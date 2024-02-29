package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
)

type Config struct {
	ConfigFilePath bool `name:"file-path" help:"Prints the path to your config file"`
	lib.NoStyleArgs
}

func (opt *Config) Help() string {
	return `You are able to configure some of klog’s behaviour by providing a configuration file in your klog config folder. (Run ` + "`" + `klog config --file-path` + "`" + ` to print the path of that config file.)

If you run ` + "`" + `klog config` + "`" + `, you can learn about the supported properties in the file, and you also see what values are in effect at the moment. (Note: the output of the command does not print the actual file, it rather shows the configuration as it is in effect!)

You may use the output as template for setting up your config file, as its format is valid as shown.`
}

func (opt *Config) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	if opt.ConfigFilePath {
		ctx.Print(app.Join(ctx.KlogConfigFolder(), app.CONFIG_FILE_NAME).Path() + "\n")
		return nil
	}
	colours := ctx.Serialiser().Colours()
	for i, e := range app.CONFIG_FILE_ENTRIES {
		ctx.Print(ctx.Serialiser().Format(colours.Subdued, lib.Reflower.Reflow(e.Help.Summary, []string{"# "})))
		ctx.Print("\n")
		ctx.Print(ctx.Serialiser().Format(colours.Subdued, lib.Reflower.Reflow("Value: "+e.Help.Value, []string{"# - ", "#   "})))
		ctx.Print("\n")
		ctx.Print(ctx.Serialiser().Format(colours.Subdued, lib.Reflower.Reflow("Default: "+e.Help.Default, []string{"# - ", "#   "})))
		ctx.Print("\n")
		ctx.Print(ctx.Serialiser().Format(colours.Red, e.Name))
		ctx.Print(" = ")
		ctx.Print(ctx.Serialiser().Format(terminalformat.Style{Color: "227"}, e.Value(ctx.Config())))
		if i < len(app.CONFIG_FILE_ENTRIES)-1 {
			ctx.Print("\n\n")
		}
	}
	ctx.Print("\n")
	return nil
}
