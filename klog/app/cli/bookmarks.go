package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"strings"
)

type Bookmarks struct {
	List BookmarksList `cmd:"" help:"Displays all bookmarks"`
	Ls   BookmarksList `cmd:"" hidden:"" help:"Alias for 'list'"`

	Set BookmarksSet `cmd:"" help:"Defines a bookmark (or overwrites an existing one)"`
	New BookmarksSet `cmd:"" hidden:"" help:"Alias for 'set'"`

	Unset BookmarksUnset `cmd:"" help:"Removes a bookmark from the collection"`
	Rm    BookmarksUnset `cmd:"" hidden:"" help:"Alias for 'unset'"`

	Clear BookmarksClear `cmd:"" help:"Clears entire bookmark collection"`

	Info BookmarksInfo `cmd:"" help:"Prints file information for a bookmark"`
}

func (opt *Bookmarks) Help() string {
	return `Bookmarks allow you to interact with often-used files via an alias,
regardless of your current working directory. A bookmark name is always prefixed with an '@'.

E.g.: klog total @work

You can specify as many bookmarks as you want. There can even be one “unnamed” bookmark.`
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
	Name string `arg:"" name:"bookmark" type:"string" predictor:"bookmark" help:"The path of the bookmark"`
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
	File  string `arg:"" type:"string" predictor:"file" help:".klg source file"`
	Name  string `arg:"" name:"bookmark" type:"string" optional:"1" help:"The name of the bookmark."`
	Force bool   `name:"force" help:"Force to set, even if target file does not exist or is invalid"`
	lib.QuietArgs
}

func (opt *BookmarksSet) Run(ctx app.Context) error {
	file, err := app.NewFile(opt.File)
	if err != nil {
		return err
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
			ctx.Print("Changed bookmark:\n")
		} else {
			ctx.Print("Created new bookmark:\n")
		}
	}
	ctx.Print(bookmark.Name().ValuePretty() + " -> " + bookmark.Target().Path() + "\n")
	return nil
}

type BookmarksUnset struct {
	// The name is not optional here, to avoid accidental invocations
	Name string `arg:"" name:"bookmark" type:"string" predictor:"bookmark" help:"The name of the bookmark"`
	lib.QuietArgs
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
	lib.QuietArgs
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
