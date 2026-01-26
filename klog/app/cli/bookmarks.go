package cli

import (
	"strings"

	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/args"
)

type Bookmarks struct {
	List BookmarksList `cmd:"" help:"Display all bookmarks."`
	Ls   BookmarksList `cmd:"" hidden:"" help:"Alias for 'list'."`

	Set BookmarksSet `cmd:"" help:"Define a bookmark (or overwrite an existing one)."`
	New BookmarksSet `cmd:"" hidden:"" help:"Alias for 'set'."`

	Unset BookmarksUnset `cmd:"" help:"Remove a bookmark from the collection. (This only removes the bookmark, not the target file.)"`
	Rm    BookmarksUnset `cmd:"" hidden:"" help:"Alias for 'unset'."`

	Clear BookmarksClear `cmd:"" help:"Clear entire bookmark collection. (This only removes the bookmarks, not the target files.)"`

	Info BookmarksInfo `cmd:"" help:"Print file information for a bookmark."`
}

func (opt *Bookmarks) Help() string {
	return `
Bookmarks allow you to interact with often-used files via an alias, independent of your current working directory.
You can think of a bookmark as some sort of klog-specific symlink, that’s always available when you invoke klog, and that resolves to the designated target file.
Use the subcommands below to set up and manage your bookmarks.

A bookmark name is denoted by the prefix '@'. For example, if you have a bookmark named '@work', that points to a .klg file, you can use klog like this:

    klog total @work
    klog start --summary 'Started new project' @work
    klog edit @work

You can specify as many bookmarks as you want. There can also be one “unnamed” default bookmark (which internally is identified by the name '@default').
This is useful in case you only have one main file at a time, and allows you to use klog without any file-related input arguments at all. E.g.:

    klog total
    klog start --summary 'Started new project'

When setting up a bookmark, you can also create the respective target file on the file system by using the '--create' flag.

Note that klog keeps track of the bookmarks in an internal config file. Run 'klog config --help' to learn more.
`
}

type BookmarksList struct{}

func (opt *BookmarksList) Run(ctx app.Context) app.Error {
	bc, err := ctx.ReadBookmarks()
	if err != nil {
		return err
	}
	if bc.Count() == 0 {
		ctx.Print("There are no bookmarks defined yet.\n")
		return nil
	}
	for _, b := range bc.All() {
		ctx.Print(b.Name().ValuePretty() + " -> " + b.Target().Path() + "\n")
	}
	return nil
}

type BookmarksInfo struct {
	Dir  bool   `name:"dir" type:"string" help:"Display the directory"`
	File bool   `name:"file" type:"string" help:"Display the file name"`
	Name string `arg:"" name:"bookmark" type:"string" completion-predictor:"bookmark" help:"The path of the bookmark"`
}

func (opt *BookmarksInfo) Run(ctx app.Context) error {
	bc, err := ctx.ReadBookmarks()
	if err != nil {
		return err
	}
	bookmark := bc.Get(app.NewName(opt.Name))
	if bookmark == nil {
		return app.NewErrorWithCode(
			app.NO_SUCH_BOOKMARK_ERROR,
			"No such bookmark",
			"There is no bookmark with that alias",
			nil,
		)
	}
	if opt.Dir {
		ctx.Print(bookmark.Target().Location() + "\n")
	} else if opt.File {
		ctx.Print(bookmark.Target().Name() + "\n")
	} else {
		ctx.Print(bookmark.Target().Path() + "\n")
	}
	return nil
}

type BookmarksSet struct {
	File   string `arg:"" type:"string" completion-predictor:"file" help:".klg target file"`
	Name   string `arg:"" name:"bookmark" type:"string" optional:"1" help:"The name of the bookmark."`
	Create bool   `name:"create" short:"c" help:"Create the target file"`
	Force  bool   `name:"force" help:"Force to set, even if target file does not exist or is invalid"`
	args.QuietArgs
}

func (opt *BookmarksSet) Run(ctx app.Context) error {
	file, err := app.NewFile(opt.File)
	if err != nil {
		return err
	}

	if opt.Create {
		cErr := app.CreateEmptyFile(file)
		if cErr != nil {
			return cErr
		}
	}

	if !opt.Force {
		_, rErr := ctx.ReadInputs(app.FileOrBookmarkName(file.Path()))
		if rErr != nil {
			return app.NewErrorWithCode(
				app.GENERAL_ERROR,
				"Invalid bookmark target",
				"Please check that the file exists and is valid",
				rErr,
			)
		}
	}
	bookmark := (func() app.Bookmark {
		if opt.Name == "" {
			return app.NewDefaultBookmark(file)
		}
		return app.NewBookmark(opt.Name, file)
	})()
	didBookmarkAlreadyExist := false
	mErr := ctx.ManipulateBookmarks(func(bc app.BookmarksCollection) app.Error {
		didBookmarkAlreadyExist = bc.Get(bookmark.Name()) != nil
		bc.Set(bookmark)
		return nil
	})
	if mErr != nil {
		return mErr
	}
	if !opt.Quiet {
		if didBookmarkAlreadyExist {
			ctx.Print("Changed bookmark")
		} else {
			ctx.Print("Created new bookmark")
		}
		if opt.Create {
			ctx.Print(" and created target file")
		}
		ctx.Print(":\n")
	}
	ctx.Print(bookmark.Name().ValuePretty() + " -> " + bookmark.Target().Path() + "\n")
	return nil
}

type BookmarksUnset struct {
	// The name is not optional here, to avoid accidental invocations
	Name string `arg:"" name:"bookmark" type:"string" completion-predictor:"bookmark" help:"The name of the bookmark"`
	args.QuietArgs
}

func (opt *BookmarksUnset) Run(ctx app.Context) error {
	name := app.NewName(opt.Name)
	err := ctx.ManipulateBookmarks(func(bc app.BookmarksCollection) app.Error {
		hasRemoved := bc.Remove(name)
		if !hasRemoved {
			return app.NewErrorWithCode(
				app.NO_SUCH_BOOKMARK_ERROR,
				"No such bookmark",
				"Name: "+name.ValuePretty(),
				nil,
			)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if !opt.Quiet {
		ctx.Print("Removed bookmark " + name.ValuePretty() + "\n")
	}
	return nil
}

type BookmarksClear struct {
	Yes bool `name:"yes" short:"y" help:"Skip confirmation"`
	args.QuietArgs
}

func (opt *BookmarksClear) Run(ctx app.Context) error {
	if !opt.Yes {
		ctx.Print("Do you want to clear all bookmarks? [y/N] ")
		confirmation, err := ctx.ReadLine()
		if err != nil {
			return err
		}
		if strings.ToLower(confirmation) != "y" {
			return nil
		}
	}
	err := ctx.ManipulateBookmarks(func(bc app.BookmarksCollection) app.Error {
		bc.Clear()
		return nil
	})
	if err != nil {
		return err
	}
	if !opt.Quiet {
		ctx.Print("Cleared all bookmarks\n")
	}
	return nil
}
