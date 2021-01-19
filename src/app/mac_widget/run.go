package mac_widget

import (
	"klog/app"
	"klog/lib/caseymrm/menuet"
	"os"
	"time"
)

var ticker = time.NewTicker(500 * time.Millisecond)

func Run() {
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "⏱",
	})
	menuet.App().Name = "klog widget"
	menuet.App().Label = "-" // not actually needed, but needs to be set
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
