package cli

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
	"strings"
)

type Bookmarks struct {
	List  BookmarksList  `cmd name:"list" help:"Display all bookmarks"`
	Ls    BookmarksList  `cmd name:"ls" hidden help:"Alias"`
	Set   BookmarksSet   `cmd name:"set" help:"Define a bookmark"`
	Unset BookmarksUnset `cmd name:"unset" help:"Remove a bookmark"`
	Clear BookmarksClear `cmd name:"clear" help:"Clear entire bookmark collection"`
}

func (opt *Bookmarks) Help() string {
	return `Bookmarks allow you to interact with often-used files via a short alias,
regardless of your current working directory.

E.g.: klog total @myfile

You can specify as many bookmarks as you want. There can even be one “unnamed” bookmark.`
}

type BookmarksList struct{}

func (opt *BookmarksList) Run(ctx app.Context) error {
	bc, err := ctx.ReadBookmarks()
	if err != nil {
		return err
	}
	for _, b := range bc.All() {
		ctx.Print(b.Name().ValuePretty() + " -> " + b.Target().Path() + "\n")
	}
	return nil
}

type BookmarksSet struct {
	File string `arg type:"existingfile" help:".klg source file"`
	Name string `arg name:"bookmark" type:"string" optional:"1" help:"The name of the bookmark"`
	lib.QuietArgs
}

func (opt *BookmarksSet) Run(ctx app.Context) error {
	return ctx.ManipulateBookmarks(func(bc app.BookmarksCollection) app.Error {
		bookmark := (func() app.Bookmark {
			if opt.Name == "" {
				return app.NewDefaultBookmark(opt.File)
			}
			return app.NewBookmark(opt.Name, opt.File)
		})()
		bc.Add(bookmark)
		if opt.Quiet {
			ctx.Print(bookmark.Name().ValuePretty() + " -> " + bookmark.Target().Path())
		} else {
			ctx.Print("Created bookmark " + bookmark.Name().ValuePretty() + " for file " + bookmark.Target().Path())
		}
		ctx.Print("\n")
		return nil
	})
}

type BookmarksUnset struct {
	// The name is not optional here, to avoid accidental invocations
	Name string `arg name:"bookmark" type:"string" help:"The name of the bookmark"`
	lib.QuietArgs
}

func (opt *BookmarksUnset) Run(ctx app.Context) error {
	return ctx.ManipulateBookmarks(func(bc app.BookmarksCollection) app.Error {
		name := app.NewName(opt.Name)
		hasRemoved := bc.Remove(name)
		if !hasRemoved {
			return app.NewErrorWithCode(
				app.NO_BOOKMARK_SET_ERROR,
				"No such bookmark",
				"Name: "+name.ValuePretty(),
				nil,
			)
		}
		if !opt.Quiet {
			ctx.Print("Removed bookmark " + name.ValuePretty() + "\n")
		}
		return nil
	})
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
	return ctx.ManipulateBookmarks(func(bc app.BookmarksCollection) app.Error {
		bc.Clear()
		if !opt.Quiet {
			ctx.Print("Cleared all bookmarks\n")
		}
		return nil
	})
}
