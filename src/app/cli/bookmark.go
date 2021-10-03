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
	bc, err := ctx.ReadBookmarks()
	if err != nil {
		return err
	}
	defaultBookmark := bc.Default()
	if defaultBookmark == nil {
		return newNoBookmarkSetError()
	}
	if !opt.Quiet {
		ctx.Print("Current bookmark: ")
	}
	ctx.Print(defaultBookmark.Target().Path + "\n")
	return nil
}

type BookmarkSet struct {
	File string `arg type:"existingfile" help:".klg source file"`
	lib.QuietArgs
}

func (args *BookmarkSet) Run(ctx app.Context) error {
	bc, err := ctx.ReadBookmarks()
	if err != nil {
		return err
	}
	bc.Add(app.NewDefaultBookmark(args.File))
	err = ctx.SaveBookmarks(bc)
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
	bc, err := ctx.ReadBookmarks()
	if err != nil {
		return err
	}
	defaultBookmark := bc.Default()
	if defaultBookmark == nil {
		return newNoBookmarkSetError()
	}
	return ctx.OpenInEditor(defaultBookmark.Target().Path)
}

type BookmarkUnset struct{}

func (args *BookmarkUnset) Run(ctx app.Context) error {
	bc, err := ctx.ReadBookmarks()
	if err != nil {
		return err
	}
	bc.Clear()
	err = ctx.SaveBookmarks(bc)
	if err != nil {
		return err
	}
	ctx.Print("Cleared bookmark\n")
	return nil
}

func newNoBookmarkSetError() error {
	return app.NewErrorWithCode(
		app.NO_BOOKMARK_SET_ERROR,
		"No bookmark set",
		"You can set a bookmark by running: klog bookmark set somefile.klg",
		nil,
	)
}
