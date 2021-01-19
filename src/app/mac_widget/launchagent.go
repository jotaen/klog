package mac_widget

import (
	"io/ioutil"
	"os"
)

var launchAgentName = "net.jotaen.klog.plist"

func createLaunchAgent(homeDir string) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	contents := `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
<plist version='1.0'>
	<dict>
		<key>Label</key>
		<string>` + launchAgentName + `</string>		

		<key>ProgramArguments</key>
		<array>
			<string>` + execPath + `</string>
			<string>widget</string>
			<string>--start</string>
		</array>

		<key>RunAtLoad</key>
		<true/>
	</dict>
</plist>

`
	// chmod=0731 is how MacOS sets it
	err = os.MkdirAll(launchAgentDir(homeDir), 0731)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(launchAgentDir(homeDir)+launchAgentName, []byte(contents), 0644)
}

func removeLaunchAgent(homeDir string) error {
	return os.Remove(launchAgentDir(homeDir) + launchAgentName)
}

func hasLaunchAgent(homeDir string) bool {
	fi, err := os.Stat(launchAgentDir(homeDir) + launchAgentName)
	return err == nil && !fi.IsDir()
}

func launchAgentDir(homeDir string) string {
	return homeDir + "/Library/LaunchAgents/"
}
