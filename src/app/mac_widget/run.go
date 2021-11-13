//go:build darwin && amd64

package mac_widget

import (
	menuet2 "github.com/jotaen/klog/lib/caseymrm/menuet"
	"github.com/jotaen/klog/src/app"
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
	ctx, err := app.NewContextFromEnv(&parser.PlainSerialiser)
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

	menuet2.App().SetMenuState(&menuet2.MenuState{
		Title: "⏱",
	})
	menuet2.App().Name = "klog widget"
	menuet2.App().Label = "-" // not actually needed, but needs to be set
	menuet2.App().Children = func() []menuet2.MenuItem {
		return render(ctx, &launchAgent)
	}

	go updateTimer()
	menuet2.App().RunApplication()
}

func updateTimer() {
	for {
		<-ticker.C
		menuet2.App().MenuChanged()
	}
}
