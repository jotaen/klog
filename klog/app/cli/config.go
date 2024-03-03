package cli

import (
	"github.com/jotaen/klog/klog/app"
	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/app/cli/util"
)

type Config struct {
	util.NoStyleArgs
}

func (opt *Config) Help() string {
	return `
You are able to configure some of klog’s behaviour by providing a configuration file.

If you run 'klog config', you can learn about the supported properties in the file, and which of those you have set.
You may use the output of that command as template for setting up your config file, as its format is valid syntax.

The configuration file is named 'config.ini' and resides in your klog config folder.
Run 'klog info config-folder' to learn where your klog config folder is located.
`
}

func (opt *Config) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	styler, _ := ctx.Serialise()
	for i, e := range app.CONFIG_FILE_ENTRIES {
		ctx.Print(styler.Props(tf.StyleProps{Color: tf.SUBDUED}).Format(util.Reflower.Reflow(e.Help.Summary, []string{"# "})))
		ctx.Print("\n")
		ctx.Print(styler.Props(tf.StyleProps{Color: tf.SUBDUED}).Format(util.Reflower.Reflow("Value: "+e.Help.Value, []string{"# - ", "#   "})))
		ctx.Print("\n")
		ctx.Print(styler.Props(tf.StyleProps{Color: tf.SUBDUED}).Format(util.Reflower.Reflow("Default: "+e.Help.Default, []string{"# - ", "#   "})))
		ctx.Print("\n")
		ctx.Print(styler.Props(tf.StyleProps{Color: tf.RED}).Format(e.Name))
		ctx.Print(" = ")
		ctx.Print(styler.Props(tf.StyleProps{Color: tf.YELLOW}).Format(e.Value(ctx.Config())))
		if i < len(app.CONFIG_FILE_ENTRIES)-1 {
			ctx.Print("\n\n")
		}
	}
	ctx.Print("\n")
	return nil
}
