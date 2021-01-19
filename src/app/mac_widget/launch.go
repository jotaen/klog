package mac_widget

import (
	"fmt"
	"klog/app"
	"klog/lib/caseymrm/menuet"
	"klog/record"
	"klog/service"
	"os"
	"time"
)

var ticker = time.NewTicker(1 * time.Second)

func Launch() {
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "⏱",
	})
	menuet.App().Name = "klog widget"
	menuet.App().Label = ""
	menuet.App().Children = func() []menuet.MenuItem {
		ctx, err := app.NewContextFromEnv()
		if err != nil {
			os.Exit(1)
		}
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

	//items = append(items, menuet.MenuItem{
	//	Text: "Settings",
	//	Children: func() []menuet.MenuItem {
	//		var items []menuet.MenuItem
	//		for _, b := range ctx.LatestFiles() {
	//			items = append(items, menuet.MenuItem{
	//				Text:  b,
	//				State: ctx.OutputFilePath() == b,
	//			})
	//		}
	//		return items
	//	},
	//})

	//			items := []menuet.MenuItem{
	//				{Text: "Open folder", Clicked: func() {
	//					app.OpenInFileBrowser(project)
	//				}},
	//			}

	items = append(items, menuet.MenuItem{
		Text:  "Start widget on login",
		State: hasLaunchAgent(ctx.HomeDir()),
		Clicked: func() {
			var err error
			if hasLaunchAgent(ctx.HomeDir()) {
				err = removeLaunchAgent(ctx.HomeDir())
			} else {
				err = createLaunchAgent(ctx.HomeDir())
			}
			if err != nil {
				fmt.Println(err)
			}
		},
	})

	return items
}

func renderRecords(ctx app.Context) []menuet.MenuItem {
	var items []menuet.MenuItem
	now := time.Now()
	nowTime, _ := record.NewTime(now.Hour(), now.Minute())
	nowDate := record.NewDateFromTime(now)
	rs, _ := service.FindFilter(ctx.BookmarkedFile(), service.Filter{
		BeforeEq: nowDate, AfterEq: nowDate,
	})
	var today record.Record
	if len(rs) == 1 {
		today = rs[0]
	}

	items = append(items, menuet.MenuItem{
		Text: "My File",
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
