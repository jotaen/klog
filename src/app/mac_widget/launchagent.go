package mac_widget

import (
	"io/ioutil"
	"os"
)

type launchAgent struct {
	name           string
	klogBinPath    string
	launchAgentDir string
	plistFilePath  string
}

// os.Executable()
func NewLaunchAgent(homeDir string, klogBinPath string) launchAgent {
	name := "net.jotaen.klog.widget"
	launchAgentDir := homeDir + "/Library/LaunchAgents/"
	return launchAgent{
		name:           name,
		klogBinPath:    klogBinPath,
		launchAgentDir: launchAgentDir,
		plistFilePath:  launchAgentDir + name + ".plist",
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
	err = ioutil.WriteFile(l.plistFilePath, []byte(contents), 0644)
	return err
}

func (l *launchAgent) deactivate() error {
	return os.Remove(l.plistFilePath)
}

func (l *launchAgent) isActive() bool {
	fi, err := os.Stat(l.plistFilePath)
	return err == nil && !fi.IsDir()
}
