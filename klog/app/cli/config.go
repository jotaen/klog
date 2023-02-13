package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
)

type Config struct{}

func (opt *Config) Help() string {
	return `By placing a YAML file at ~/` + app.KLOG_FOLDER + app.CONFIG_FILE + ` you are able to configure some of klog’s behaviour.

If you run ` + "`" + `klog config` + "`" + `, you can learn about the supported YAML properties in the file, and you also see what values are in effect at the moment. (Note: the output of the command does not print the actual file!)`
}

func (opt *Config) Run(ctx app.Context) app.Error {
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
