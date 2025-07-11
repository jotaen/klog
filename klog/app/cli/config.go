package cli

import (
	"github.com/jotaen/klog/klog/app"
	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/app/cli/util"
	"path/filepath"
	"strings"
)

type Config struct {
	util.NoStyleArgs
	Location bool `name:"location" help:"Print the location of the config folder."`
}

func (opt *Config) Help() string {
	lookupOrder := func() string {
		lookups := make([]string, len(app.KLOG_CONFIG_FOLDER))
		for i, f := range app.KLOG_CONFIG_FOLDER {
			lookups[i] = filepath.Join(f.EnvVarSymbol(), f.Location)
		}
		return strings.Join(lookups, "  ->  ")
	}()

	return `
klog relies on file-based configuration to customise some of its default behaviour and to keep track of its internal state.

Run 'klog config --location' to print the path of the folder where klog looks for the configuration.
The config folder can contain one or both of the following files:
  - '` + app.CONFIG_FILE_NAME + `': you can create this file manually to override some of klog’s default behaviour. You may use the output of the 'klog config' command as template for setting up this file, as its output is in valid syntax. 
  - '` + app.BOOKMARKS_FILE_NAME + `': if you use the bookmarks functionality, then klog uses this file as database. You are not supposed to edit this file by hand! Instead, use the 'klog bookmarks' command to manage your bookmarks.

You can customise the location of the config folder via environment variables. klog uses the following lookup precedence:
  ` + lookupOrder + `
`
}

func (opt *Config) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)

	if opt.Location {
		ctx.Print(ctx.KlogConfigFolder().Path() + "\n")
		return nil
	}

	styler, _ := ctx.Serialise()
	for i, e := range app.CONFIG_FILE_ENTRIES {
		ctx.Print(styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED}).Format(util.Reflower.Reflow(e.Help.Summary, []string{"# "})))
		ctx.Print("\n")
		ctx.Print(styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED}).Format(util.Reflower.Reflow("Value: "+e.Help.Value, []string{"# - ", "#   "})))
		ctx.Print("\n")
		ctx.Print(styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED}).Format(util.Reflower.Reflow("Default: "+e.Help.Default, []string{"# - ", "#   "})))
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
