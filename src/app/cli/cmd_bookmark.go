package cli

import (
	"klog/app"
)

type Bookmark struct {
	Get   BookmarkGet   `cmd name:"get" group:"Bookmark" help:"Show current bookmark"`
	Set   BookmarkSet   `cmd name:"set" group:"Bookmark" help:"Set bookmark to a file"`
	Edit  BookmarkEdit  `cmd name:"edit" group:"Bookmark" help:"Open bookmark in your editor"`
	Unset BookmarkUnset `cmd name:"unset" group:"Bookmark" help:"Clear current bookmark"`
}

type BookmarkGet struct{}

func (opt *BookmarkGet) Run(ctx app.Context) error {
	b, err := ctx.Bookmark()
	if err != nil {
		return err
	}
	ctx.Print("Current bookmark: " + b.Path + "\n")
	return nil
}

type BookmarkSet struct {
	File string `arg type:"existingfile" help:".klg source file"`
}

func (args *BookmarkSet) Run(ctx app.Context) error {
	err := ctx.SetBookmark(args.File)
	if err != nil {
		return err
	}
	ctx.Print("Bookmarked file " + args.File + "\n")
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
