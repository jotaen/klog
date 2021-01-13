package tray

import (
	"github.com/caseymrm/menuet"
	"klog/app"
	"klog/project"
	"time"
)

var currentProject project.Project
var ticker = time.NewTicker(1 * time.Second)

func Start(env app.Environment) {
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "⏱",
	})
	menuet.App().Name = "klog menu bar app"
	menuet.App().Label = "net.jotaen.klog.menuapp"
	menuet.App().Children = func() []menuet.MenuItem { return render(env) }

	currentProject = env.SavedProjects()[0]

	go updateTimer()
	menuet.App().RunApplication()
}

func updateTimer() {
	for {
		<-ticker.C
		menuet.App().MenuChanged()
	}
}

func render(env app.Environment) []menuet.MenuItem {
	var items []menuet.MenuItem

	if currentProject != nil {
		items = append(items, renderProject(currentProject)...)
		items = append(items, menuet.MenuItem{
			Type: menuet.Separator,
		})
	}

	items = append(items, menuet.MenuItem{
		Text: "Projects",
		Children: func() []menuet.MenuItem {
			var items []menuet.MenuItem
			for _, p := range env.SavedProjects() {
				items = append(items, menuet.MenuItem{
					Text:  p.Name(),
					State: currentProject == p,
				})
			}
			return items
		},
	})

	return items
}
