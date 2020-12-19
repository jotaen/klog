package tray

import (
	"fmt"
	"github.com/getlantern/systray"
	"os"
	"time"
)

func Start() {
	systray.Quit()
	systray.Run(onReady, func() {})
}

func onReady() {
	systray.SetTitle("⏱")

	workDays := systray.AddMenuItem("Example Inc.", "")
	workDays.AddSubMenuItem("2020-12-17", "")
	workDays.AddSubMenuItem("2020-12-16", "")
	workDays.AddSubMenuItem("2020-12-13", "")
	workDays.AddSubMenuItem("2020-12-12", "")
	workDays.AddSubMenuItem("2020-12-11", "")
	workDays.AddSubMenuItem("2020-12-03", "")
	workDays.AddSubMenuItem("2020-12-01", "")

	clock := systray.AddMenuItem("", "")
	clock.Disable()
	systray.AddMenuItem("Stop time", "")

	systray.AddSeparator()
	projects := systray.AddMenuItem("Projects", "")
	projects.AddSubMenuItemCheckbox("Example Inc.", "~/Clients/example", true)
	projects.AddSubMenuItemCheckbox("Sample Ltd.", "~/Clients/sample", false)
	projects.AddSubMenuItemCheckbox("Test & Sons", "~/Clients/test", false)
	config := systray.AddMenuItem("klog", "")
	config.AddSubMenuItemCheckbox("Launch tray app on login", "", true)
	mQuit := config.AddSubMenuItem("Quit tray app", "")

	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				hours := fmt.Sprintf("%dh %dm", now.Hour(), now.Minute())
				liveIndicator := " ⏱"
				clock.SetTitle("Today: " + hours + liveIndicator)
			}

		}
	}()

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}
