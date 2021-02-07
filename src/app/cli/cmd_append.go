package cli

import (
	"klog/app"
)

type Append struct {
	File string `arg optional type:"existingfile" name:"file" help:".klg source file (if empty the bookmark is used)"`
	From string `required name:"from" help:"The name of the template to instantiate"`
}

func (args *Append) Run(ctx app.Context) error {
	target := args.File
	if target == "" {
		b, err := ctx.Bookmark()
		if err != nil {
			return err
		}
		target = b.Path
	}
	return ctx.AppendTemplateToFile(target, args.From)
}
