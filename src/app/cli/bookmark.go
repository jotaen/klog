package cli

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
)

type Bookmark struct {
	Get   BookmarkGet   `cmd name:"get" group:"Bookmark" help:"Show current bookmark"`
	Set   BookmarkSet   `cmd name:"set" group:"Bookmark" help:"Set bookmark to a file"`
	Edit  BookmarkEdit  `cmd name:"edit" group:"Bookmark" help:"Open bookmark in your editor"`
	Unset BookmarkUnset `cmd name:"unset" group:"Bookmark" help:"Clear current bookmark"`
}

func (opt *Bookmark) Help() string {
	return `With bookmarks you can make klog always read from a default file, in case you don’t specify one explicitly.

This is handy in case you always use the same file.
You can then interact with it regardless of your current working directory.`
}

type BookmarkGet struct {
	lib.QuietArgs
}

func (opt *BookmarkGet) Run(ctx app.Context) error {
	b, err := ctx.Bookmark()
	if err != nil {
		return err
	}
	if !opt.Quiet {
		ctx.Print("Current bookmark: ")
	}
	ctx.Print(b.Path + "\n")
	return nil
}

type BookmarkSet struct {
	File string `arg type:"existingfile" help:".klg source file"`
	lib.QuietArgs
}

func (args *BookmarkSet) Run(ctx app.Context) error {
	err := ctx.SetBookmark(args.File)
	if err != nil {
		return err
	}
	if !args.Quiet {
		ctx.Print("Bookmarked file ")
	}
	ctx.Print(args.File + "\n")
	return nil
}

type BookmarkEdit struct{}

func (args *BookmarkEdit) Run(ctx app.Context) error {
	b, appErr := ctx.Bookmark()
	if appErr != nil {
		return appErr
	}
	return ctx.OpenInEditor(b.Path)
}

type BookmarkUnset struct{}

func (args *BookmarkUnset) Run(ctx app.Context) error {
	err := ctx.UnsetBookmark()
	if err != nil {
		return err
	}
	ctx.Print("Cleared bookmark\n")
	return nil
}
