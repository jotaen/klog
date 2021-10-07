package cli

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
	"strings"
)

type Bookmarks struct {
	List BookmarksList `cmd name:"list" help:"Displays all bookmarks"`
	Ls   BookmarksList `cmd name:"ls" hidden help:"Alias for 'list'"`

	Set BookmarksSet `cmd name:"set" help:"Defines a bookmark (or overwrites an existing one)"`
	New BookmarksSet `cmd name:"new" hidden help:"Alias for 'set'"`

	Unset BookmarksUnset `cmd name:"unset" help:"Removes a bookmark from the collection"`
	Rm    BookmarksUnset `cmd name:"rm" hidden help:"Alias for 'unset'"`

	Clear BookmarksClear `cmd name:"clear" help:"Clears entire bookmark collection"`
}

func (opt *Bookmarks) Help() string {
	return `Bookmarks allow you to interact with often-used files via an alias,
regardless of your current working directory. A bookmark name is always prefixed with an '@'.

E.g.: klog total @work

You can specify as many bookmarks as you want. There can even be one “unnamed” bookmark.`
}

type BookmarksList struct{}

func (opt *BookmarksList) Run(ctx app.Context) error {
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

type BookmarksSet struct {
	File  string `arg type:"string" help:".klg source file"`
	Name  string `arg name:"bookmark" type:"string" optional:"1" help:"The name of the bookmark."`
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
	mErr := ctx.ManipulateBookmarks(func(bc app.BookmarksCollection) app.Error {
		bc.Set(bookmark)
		return nil
	})
	if mErr != nil {
		return mErr
	}
	if !opt.Quiet {
		ctx.Print("Created new bookmark:\n")
	}
	ctx.Print(bookmark.Name().ValuePretty() + " -> " + bookmark.Target().Path() + "\n")
	return nil
}

type BookmarksUnset struct {
	// The name is not optional here, to avoid accidental invocations
	Name string `arg name:"bookmark" type:"string" help:"The name of the bookmark"`
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
