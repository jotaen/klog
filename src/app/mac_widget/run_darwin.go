package mac_widget

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/lib/caseymrm/menuet"
	"github.com/jotaen/klog/src/parser"
	"os"
	"os/exec"
	"time"
)

var ticker = time.NewTicker(500 * time.Millisecond)

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
		_ = exec.Command("launchctl", "load", launchAgent.plistFilePath).Run()
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
