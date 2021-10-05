package mac_widget

import (
	"github.com/jotaen/klog/src/app"
	"os"
)

type launchAgent struct {
	name           string
	klogBinPath    string
	launchAgentDir string
	plistFile      app.File
}

func newLaunchAgent(homeDir string, klogBinPath string) launchAgent {
	name := "net.jotaen.klog.widget"
	launchAgentDir := homeDir + "/Library/LaunchAgents/"
	return launchAgent{
		name:           name,
		klogBinPath:    klogBinPath,
		launchAgentDir: launchAgentDir,
		plistFile:      app.NewFileOrPanic(launchAgentDir + name + ".plist"),
	}
}

func (l *launchAgent) activate() error {
	contents := `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
<plist version='1.0'>
	<dict>
		<key>Label</key>
		<string>` + l.name + `</string>		

		<key>ProgramArguments</key>
		<array>
			<string>` + l.klogBinPath + `</string>
			<string>widget</string>
		</array>

		<key>RunAtLoad</key>
		<true/>
	</dict>
</plist>

`
	// chmod=0731 is how MacOS sets it
	err := os.MkdirAll(l.launchAgentDir, 0731)
	if err != nil {
		return err
	}
	err = app.WriteToFile(l.plistFile, contents)
	return err
}

func (l *launchAgent) deactivate() error {
	return os.Remove(l.plistFile.Path())
}

func (l *launchAgent) isActive() bool {
	fi, err := os.Stat(l.plistFile.Path())
	return err == nil && !fi.IsDir()
}
