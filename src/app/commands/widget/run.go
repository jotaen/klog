//go:build darwin && amd64

/*
Package widget is a native Mac Widget. It’s not maintained anymore and might
be discontinued.
*/
package widget

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/commands/widget/menuet"
	"github.com/jotaen/klog/src/parser"
	"os"
	"os/exec"
	"time"
)

var ticker = time.NewTicker(500 * time.Millisecond)

func IsWidgetAvailable() bool {
	return true
}

func Run(forceRunThroughLaunchAgent bool) {
	ctx, err := app.NewContextFromEnv(app.Meta{}, &parser.PlainSerialiser)
	if err != nil {
		os.Exit(1)
	}
	binPath, _ := os.Executable()
	launchAgent := newLaunchAgent(ctx.HomeFolder(), binPath)

	if forceRunThroughLaunchAgent {
		if !launchAgent.isActive() {
			_ = launchAgent.activate()
		}
		_ = exec.Command("launchctl", "load", launchAgent.plistFile.Path()).Run()
		_ = exec.Command("launchctl", "start", launchAgent.name).Run()
		os.Exit(0)
	}

	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "⏱",
	})
	menuet.App().Name = "klog widget"
	menuet.App().Label = "-" // not actually needed, but needs to be set
	menuet.App().Children = func() []menuet.MenuItem {
		return render(ctx, &launchAgent)
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
