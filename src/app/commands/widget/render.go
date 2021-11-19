//go:build darwin && amd64

package widget

import (
	"fmt"
	menuet2 "github.com/jotaen/klog/lib/caseymrm/menuet"
	"github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/service"
)

var blinker = blinkerT{1}

func render(ctx app.Context, agent *launchAgent) []menuet2.MenuItem {
	var items []menuet2.MenuItem

	items = append(items, func() []menuet2.MenuItem {
		bc, err := ctx.ReadBookmarks()
		if err != nil || bc.Default() == nil {
			return []menuet2.MenuItem{{
				Text:       "No bookmark specified",
				FontWeight: menuet2.WeightBold,
			}, {
				Text: "Bookmark a file by running:",
			}, {
				Text: "klog bookmark set yourfile.klg",
			}}
		}
		defaultBookmark := bc.Default()
		rs, pErr := ctx.ReadInputs()
		if pErr != nil {
			return []menuet2.MenuItem{{
				Text: defaultBookmark.Target().Name(),
			}, {
				Text: "Error: file cannot be parsed",
			}}
		}
		return renderRecords(ctx, rs, defaultBookmark.Target())
	}()...)

	items = append(items, menuet2.MenuItem{
		Type: menuet2.Separator,
	}, menuet2.MenuItem{
		Text: "Settings",
		Children: func() []menuet2.MenuItem {
			return []menuet2.MenuItem{{
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

func renderRecords(ctx app.Context, records []klog.Record, file app.File) []menuet2.MenuItem {
	var items []menuet2.MenuItem

	today := service.Filter(records, service.FilterQry{Dates: []klog.Date{klog.NewDateFromTime(ctx.Now())}})
	if today != nil {
		total, isOngoing := service.HypotheticalTotal(ctx.Now(), today...)
		indicator := ""
		if isOngoing {
			indicator = "  " + blinker.blink()
		}
		items = append(items, menuet2.MenuItem{
			Text: "Today: " + total.ToString() + indicator,
		})
	}

	items = append(items, menuet2.MenuItem{
		Text: file.Name(),
		Children: func() []menuet2.MenuItem {
			total := service.Total(records...)
			should := service.ShouldTotalSum(records...)
			diff := service.Diff(should, total)
			plus := ""
			if diff.InMinutes() > 0 {
				plus = "+"
			}
			items := []menuet2.MenuItem{
				{
					Text: "Show in Finder...",
					Clicked: func() {
						_ = ctx.OpenInFileBrowser(file)
					},
				},
				{Type: menuet2.Separator},
				{Text: "Total: " + total.ToString()},
				{Text: "Should: " + should.ToString()},
				{Text: "Diff: " + plus + diff.ToString()},
			}
			showMax := 7
			for i, r := range service.Sort(records, false) {
				if i == 0 {
					items = append(items, menuet2.MenuItem{Type: menuet2.Separator})
				}
				if i == showMax {
					items = append(items, menuet2.MenuItem{Text: fmt.Sprintf("(%d more)", len(records)-showMax)})
					break
				}
				items = append(items, menuet2.MenuItem{Text: r.Date().ToString() + ": " + service.Total(r).ToString()})
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
