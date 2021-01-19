package mac_widget

import (
	"klog/app"
	"os"
)

func createLaunchAgent(homeDir string) error {
	content := `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
<plist version='1.0'>
	<dict>
		<key>Label</key>
  		<string>net.jotaen.klog.plist</string>

		<key>Program</key>
		<string>klog</string>

		<key>ProgramArguments</key>
		<array>
			<string>widget</string>
			<string>--start</string>
		</array>
     
		<key>RunAtLoad</key>
		<true/>
	</dict>
</plist>
`
	return app.WriteFile(launchAgentPath(homeDir), content)
}

func removeLaunchAgent(homeDir string) error {
	return os.Remove(launchAgentPath(homeDir))
}

func hasLaunchAgent(homeDir string) bool {
	fi, err := os.Stat(launchAgentPath(homeDir))
	return err == nil && !fi.IsDir()
}

func launchAgentPath(homeDir string) string {
	return homeDir + "/Library/LaunchAgents/net.jotaen.klog.plist"
}
