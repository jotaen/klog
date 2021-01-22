package mac_widget

import (
	"fmt"
	"klog"
	"klog/app"
	"klog/lib/caseymrm/menuet"
	"klog/service"
	"os/exec"
)

var blinker = blinkerT{1}

func render(ctx *app.Context, agent *launchAgent) []menuet.MenuItem {
	var items []menuet.MenuItem

	items = append(items, func() []menuet.MenuItem {
		if rs, file, err := ctx.Bookmark(); err == nil {
			return renderRecords(rs, file)
		}
		return []menuet.MenuItem{{
			Text:       "No input file specified",
			FontWeight: menuet.WeightBold,
		}, {
			Text: "Specify one by running:",
		}, {
			Text: "klog widget --file yourfile.klg",
		}}
	}()...)

	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
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

func renderRecords(records []src.Record, file app.File) []menuet.MenuItem {
	var items []menuet.MenuItem

	items = append(items, menuet.MenuItem{
		Text: file.Name,
		Children: func() []menuet.MenuItem {
			total := service.Total(records...)
			should := service.ShouldTotal(records...)
			diff := src.NewDuration(0, 0).Minus(should).Plus(total)
			plus := ""
			if diff.InMinutes() > 0 {
				plus = "+"
			}
			items := []menuet.MenuItem{
				{
					Text: "Show in Finder...",
					Clicked: func() {
						cmd := exec.Command("open", file.Location)
						_ = cmd.Start()
					},
				},
				{Type: menuet.Separator},
				{Text: "Total: " + total.ToString()},
				{Text: "Should: " + should.ToString() + "!"},
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

	//items = append(items, menuet.MenuItem{
	//	Text: func() string {
	//		if isRunningCurrently {
	//			return "Stop"
	//		}
	//		return "Start tracking"
	//	}(),
	//	Clicked: func() {
	//		if isRunningCurrently {
	//			// stop!
	//		} else {
	//			// start!
	//		}
	//	},
	//})

	//if today != nil {
	//	items = append(items, menuet.MenuItem{
	//		State: isRunningCurrently,
	//		Text:  "Today: " + totalTimeValue,
	//	})
	//}

	return items
}

type blinkerT struct {
	cycle int
}

func (b *blinkerT) get() string {
	b.cycle++
	switch b.cycle {
	case 1:
		return "◷"
	case 2:
		return "◶"
	case 3:
		return "◵"
	}
	b.cycle = 1
	return "◴"
}
