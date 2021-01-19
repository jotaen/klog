package menuet

import (
	"log"
	"strings"
	"time"

	"github.com/caseymrm/askm"
)

// ItemType represents what type of menu item this is
type ItemType string

const (
	// Regular is a normal item with text and optional callback
	Regular ItemType = ""
	// Separator is a horizontal line
	Separator = "separator"
	// Root is the top level menu directly off the menubar
	Root = "root"
	// TODO: StartAtLogin, Quit, Image, Spinner, etc
)

// MenuItem represents one item in the dropdown
type MenuItem struct {
	Type ItemType

	Text       string
	Image      string // In Resources dir or URL, should have height 16
	FontSize   int    // Default: 14
	FontWeight FontWeight
	State      bool // shows checkmark when set

	Clicked  func()            `json:"-"`
	Children func() []MenuItem `json:"-"`
}

type internalItem struct {
	Unique       string
	ParentUnique string
	HasChildren  bool
	Clickable    bool

	MenuItem
}

func (a *Application) children(unique string) []internalItem {
	a.visibleMenuItemsMutex.RLock()
	item, ok := a.visibleMenuItems[unique]
	a.visibleMenuItemsMutex.RUnlock()
	if strings.HasSuffix(unique, ":root") {
		// Create synthetic item
		item.Unique = unique
		item.Type = Root
		item.Children = a.Children
		ok = true
	}
	if !ok {
		log.Printf("Item not found for children: %s", unique)
	}
	var items []MenuItem
	if item.Children != nil {
		items = item.Children()
	}
	internalItems := make([]internalItem, len(items))
	for ind, item := range items {
		a.visibleMenuItemsMutex.Lock()
		newUnique := askm.ArbitraryKeyNotInMap(a.visibleMenuItems)
		internal := internalItem{
			Unique:       newUnique,
			ParentUnique: unique,
			MenuItem:     item,
		}
		if internal.Children != nil {
			internal.HasChildren = true
		}
		if internal.Clicked != nil {
			internal.Clickable = true
		}
		a.visibleMenuItems[newUnique] = internal
		internalItems[ind] = internal
		a.visibleMenuItemsMutex.Unlock()
	}
	return internalItems
}

func (a *Application) menuClosed(unique string) {
	go func() {
		// We receive menuClosed before clicked, so wait a second before discarding the data just in case
		time.Sleep(100 * time.Millisecond)
		a.visibleMenuItemsMutex.Lock()
		for itemUnique, item := range a.visibleMenuItems {
			if item.ParentUnique == unique {
				delete(a.visibleMenuItems, itemUnique)
			}
		}
		a.visibleMenuItemsMutex.Unlock()
	}()
}
