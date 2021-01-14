package mac_widget

import (
	"github.com/caseymrm/menuet"
	"klog/app"
	"klog/record"
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
		service := app.NewService(nil)
		return render(service)
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

func render(service app.Service) []menuet.MenuItem {
	var items []menuet.MenuItem

	if service.Input() != nil {
		items = append(items, renderRecords(service)...)
		items = append(items, menuet.MenuItem{
			Type: menuet.Separator,
		})
	}

	items = append(items, menuet.MenuItem{
		Text: "Bookmarks",
		Children: func() []menuet.MenuItem {
			var items []menuet.MenuItem
			for _, b := range service.BookmarkedFiles() {
				items = append(items, menuet.MenuItem{
					Text:  b,
					State: service.CurrentFile() == b,
				})
			}
			return items
		},
	})

	return items
}

func renderRecords(service app.Service) []menuet.MenuItem {
	var items []menuet.MenuItem
	now := time.Now()
	nowTime, _ := record.NewTime(now.Hour(), now.Minute())
	nowDate, _ := record.NewDateFromTime(now)
	today := record.Find(nowDate, service.Input())

	items = append(items, menuet.MenuItem{
		Text: service.CurrentFile(),
	})

	totalTimeValue := (func() string {
		if today != nil {
			if today.OpenRange() != nil {
				untilNow, _ := record.NewRange(today.OpenRange(), nowTime)
				if untilNow != nil {
					result := ""
					result = record.Total(today).Add(untilNow.Duration()).ToString()
					result += "  "
					if now.Second()%2 == 0 {
						result += "◑"
					} else if now.Second()%2 == 1 {
						result += "◐"
					}
					return result
				}
			} else {
				return record.Total(today).ToString()
			}
		}
		return "–"
	})()

	isRunning := today != nil && today.OpenRange() != nil
	items = append(items, menuet.MenuItem{
		Text: "Today: " + totalTimeValue,
	}, menuet.MenuItem{
		Text:  "Run Timer",
		State: isRunning,
		Clicked: func() {
			if isRunning {
				service.QuickStopAt(nowDate, nowTime)
			} else {
				service.QuickStartAt(nowDate, nowTime)
			}
		},
	})

	return items
}
