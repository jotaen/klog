package tray

import (
	"github.com/caseymrm/menuet"
	"klog/app"
	"klog/datetime"
	"klog/project"
	"time"
)

func renderProject(project project.Project) []menuet.MenuItem {
	var items []menuet.MenuItem

	items = append(items, menuet.MenuItem{
		Text: project.Name(),
	})

	now := time.Now()
	nowTime, _ := datetime.NewTime(now.Hour(), now.Minute())
	nowDate, _ := datetime.NewDateFromTime(now)
	currentDay, _ := project.Get(nowDate)
	totalTimeValue := "–"

	if currentDay != nil {
		if currentDay.OpenRange() != nil {
			untilNow, _ := datetime.NewTimeRange(currentDay.OpenRange(), nowTime)
			if untilNow != nil {
				totalTimeValue = currentDay.TotalWorkTime().Add(untilNow.Duration()).ToString()
				totalTimeValue += "  "
				if now.Second()%2 == 0 {
					totalTimeValue += "◑"
				} else if now.Second()%2 == 1 {
					totalTimeValue += "◐"
				}
			}
		} else {
			totalTimeValue = currentDay.TotalWorkTime().ToString()
		}
	}

	isRunning := currentDay != nil && currentDay.OpenRange() != nil
	items = append(items, menuet.MenuItem{
		Text: "Today: " + totalTimeValue,
	}, menuet.MenuItem{
		Text:  "Run Timer",
		State: isRunning,
		Clicked: func() {
			now := time.Now()
			if isRunning {
				app.Stop(project, now)
			} else {
				app.Start(project, now)
			}
		},
	})

	items = append(items, menuet.MenuItem{
		Text: "History",
		Children: func() []menuet.MenuItem {
			items := []menuet.MenuItem{
				{Text: "Open folder", Clicked: func() {
					app.OpenInFileBrowser(project)
				}},
			}
			days, _ := project.List()
			if len(days) > 0 {
				items = append(items, menuet.MenuItem{Type: menuet.Separator})
				for i, d := range days {
					if i == 5 {
						break
					}
					wd, _ := project.Get(d)
					if wd == nil {
						break
					}
					todayLabel := ""
					if currentDay != nil && wd.Date() == currentDay.Date() {
						todayLabel = " (today)"
					}
					items = append(items, menuet.MenuItem{
						Text: wd.Date().ToString() + todayLabel, Clicked: func() {
							app.OpenInEditor(project, wd)
						},
					})
				}
			}
			return items
		},
	})

	return items
}
