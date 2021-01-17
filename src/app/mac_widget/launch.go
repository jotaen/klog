package mac_widget

import (
	"github.com/caseymrm/menuet"
	"klog/app"
	"klog/record"
	"klog/service"
	"time"
)

var ticker = time.NewTicker(1 * time.Second)

func Launch() {
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "⏱",
	})
	menuet.App().Name = "klog widget"
	menuet.App().Label = "net.jotaen.klog.widget"
	menuet.App().Children = func() []menuet.MenuItem {
		ctx, _ := app.NewContextFromEnv() // TODO error handling
		return render(*ctx)
	}

	go updateTimer()
	menuet.App().RunApplication()
}

func updateTimer() {
	for {
		<-ticker.C
		menuet.App().MenuChanged()
	}
}

func render(ctx app.Context) []menuet.MenuItem {
	var items []menuet.MenuItem

	if ctx.BookmarkedFile() != nil {
		items = append(items, renderRecords(ctx)...)
		items = append(items, menuet.MenuItem{
			Type: menuet.Separator,
		})
	}

	items = append(items, menuet.MenuItem{
		Text: "File History",
		Children: func() []menuet.MenuItem {
			var items []menuet.MenuItem
			for _, b := range ctx.LatestFiles() {
				items = append(items, menuet.MenuItem{
					Text:  b,
					State: ctx.OutputFilePath() == b,
				})
			}
			return items
		},
	})

	return items
}

func renderRecords(ctx app.Context) []menuet.MenuItem {
	var items []menuet.MenuItem
	now := time.Now()
	nowTime, _ := record.NewTime(now.Hour(), now.Minute())
	nowDate, _ := record.NewDateFromTime(now)
	today := service.Find(nowDate, ctx.BookmarkedFile())

	items = append(items, menuet.MenuItem{
		Text: ctx.OutputFilePath(),
	})

	totalTimeValue := func() string {
		if today != nil {
			if today.OpenRange() != nil {
				untilNow, _ := record.NewRange(today.OpenRange().Start(), nowTime)
				if untilNow != nil {
					result := ""
					result = service.Total(today).Add(untilNow.Duration()).ToString()
					result += "  "
					if now.Second()%2 == 0 {
						result += "◑"
					} else if now.Second()%2 == 1 {
						result += "◐"
					}
					return result
				}
			} else {
				return service.Total(today).ToString()
			}
		}
		return "–"
	}()

	isRunning := today != nil && today.OpenRange() != nil
	items = append(items, menuet.MenuItem{
		Text: "Today: " + totalTimeValue,
	}, menuet.MenuItem{
		Text:  "Run Timer",
		State: isRunning,
		Clicked: func() {
			if isRunning {
				//ctx.QuickStopAt(nowDate, nowTime)
			} else {
				//ctx.QuickStartAt(nowDate, nowTime)
			}
		},
	})

	return items
}
