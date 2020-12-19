package tray

import (
	"github.com/caseymrm/menuet"
	"klog/app"
	"klog/datetime"
	"klog/project"
	"time"
)

var currentProject project.Project
var ticker = time.NewTicker(1 * time.Second)
var config = app.NewConfiguration("~")

func Start() {
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "⏱",
	})
	menuet.App().Name = "klog menu bar app"
	menuet.App().Label = "net.jotaen.klog.menuapp"
	menuet.App().Children = render

	currentProject = config.Projects()[0]

	go updateTimer()
	menuet.App().RunApplication()
}

func updateTimer() {
	for {
		<-ticker.C
		menuet.App().MenuChanged()
	}
}

func render() []menuet.MenuItem {
	var items []menuet.MenuItem

	if currentProject != nil {
		items = append(items, menuet.MenuItem{
			Text: currentProject.Name(),
		})

		now := time.Now()
		nowTime, _ := datetime.NewTime(now.Hour(), now.Minute())
		nowDate, _ := datetime.NewDateFromTime(now)
		currentDay, _ := currentProject.Get(nowDate)

		if currentDay != nil {
			timer, err := currentDay.TotalWorkTimeWithOpenRange(nowTime)
			if err == nil {
				stopwatch := menuet.MenuItem{
					Text: "Today: " + timer.ToString(),
				}
				items = append(items, stopwatch)
				items = append(items, menuet.MenuItem{
					Text: "Stop time",
					Clicked: func() {
						menuet.App().Alert(menuet.Alert{
							MessageText:     "Error",
							InformativeText: "Could not read file",
							Buttons:         []string{"Okay"},
						})
					},
				})
			} else {
				items = append(items, menuet.MenuItem{
					Text: "Start time",
					Clicked: func() {},
				})
			}
		}
	}

	items = append(items, menuet.MenuItem{
		Text: "History",
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{
				{Text: "Create", Clicked: func() {}},
				{Text: "Open folder", Clicked: func() {}},
				{Type: menuet.Separator},
				{Text: "Latest 10", FontSize: 10},
				{Text: "2020-12-15", Clicked: func() {}},
				{Text: "2020-12-14", Clicked: func() {}},
				{Text: "2020-12-12", Clicked: func() {}},
				{Text: "2020-12-09", Clicked: func() {}},
			}
		},
	})
	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	})

	// Switch Project
	items = append(items, menuet.MenuItem{
		Text: "Switch Project",
		Children: func() []menuet.MenuItem {
			var projects []menuet.MenuItem
			for _, p := range config.Projects() {
				projects = append(projects, menuet.MenuItem{
					Text:  p.Name(),
					State: currentProject == p,
				})
			}
			return projects
		},
	})

	return items
}
