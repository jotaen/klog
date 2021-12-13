//go:build darwin && amd64

package widget

import (
	"fmt"
	"github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/commands/widget/menuet"
	"github.com/jotaen/klog/src/service"
)

var blinker = blinkerT{1}

func render(ctx app.Context, agent *launchAgent) []menuet.MenuItem {
	var items []menuet.MenuItem

	items = append(items, func() []menuet.MenuItem {
		bc, err := ctx.ReadBookmarks()
		if err != nil || bc.Default() == nil {
			return []menuet.MenuItem{{
				Text:       "No bookmark specified",
				FontWeight: menuet.WeightBold,
			}, {
				Text: "Bookmark a file by running:",
			}, {
				Text: "klog bookmark set yourfile.klg",
			}}
		}
		defaultBookmark := bc.Default()
		rs, pErr := ctx.ReadInputs()
		if pErr != nil {
			return []menuet.MenuItem{{
				Text: defaultBookmark.Target().Name(),
			}, {
				Text: "Error: file cannot be parsed",
			}}
		}
		return renderRecords(ctx, rs, defaultBookmark.Target())
	}()...)

	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	}, menuet.MenuItem{
		Text: "Note: this widget is deprecated\n" +
			"and will be removed in one\n" +
			"of the next releases of klog.",
		FontSize: 12,
	}, menuet.MenuItem{
		Text: "Settings",
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{{
				Text:  "Launch widget on login",
				State: agent.isActive(),
				Clicked: func() {
					var err error
					if agent.isActive() {
						err = agent.deactivate()
					} else {
						err = agent.activate()
					}
					if err != nil {
						fmt.Println(err)
					}
				},
			}}
		},
	})

	return items
}

func renderRecords(ctx app.Context, records []klog.Record, file app.File) []menuet.MenuItem {
	var items []menuet.MenuItem

	today := service.Filter(records, service.FilterQry{Dates: []klog.Date{klog.NewDateFromTime(ctx.Now())}})
	if today != nil {
		total, isOngoing := service.HypotheticalTotal(ctx.Now(), today...)
		indicator := ""
		if isOngoing {
			indicator = "  " + blinker.blink()
		}
		items = append(items, menuet.MenuItem{
			Text: "Today: " + total.ToString() + indicator,
		})
	}

	items = append(items, menuet.MenuItem{
		Text: file.Name(),
		Children: func() []menuet.MenuItem {
			total := service.Total(records...)
			should := service.ShouldTotalSum(records...)
			diff := service.Diff(should, total)
			plus := ""
			if diff.InMinutes() > 0 {
				plus = "+"
			}
			items := []menuet.MenuItem{
				{
					Text: "Show in Finder...",
					Clicked: func() {
						_ = ctx.OpenInFileBrowser(app.FileOrBookmarkName(file.Path()))
					},
				},
				{Type: menuet.Separator},
				{Text: "Total: " + total.ToString()},
				{Text: "Should: " + should.ToString()},
				{Text: "Diff: " + plus + diff.ToString()},
			}
			showMax := 7
			for i, r := range service.Sort(records, false) {
				if i == 0 {
					items = append(items, menuet.MenuItem{Type: menuet.Separator})
				}
				if i == showMax {
					items = append(items, menuet.MenuItem{Text: fmt.Sprintf("(%d more)", len(records)-showMax)})
					break
				}
				items = append(items, menuet.MenuItem{Text: r.Date().ToString() + ": " + service.Total(r).ToString()})
			}
			return items
		},
	})

	return items
}

type blinkerT struct {
	cycle int
}

func (b *blinkerT) blink() string {
	b.cycle++
	switch b.cycle {
	case 1:
		return "◷"
	case 2:
		return "◶"
	case 3:
		return "◵"
	}
	b.cycle = 0
	return "◴"
}
