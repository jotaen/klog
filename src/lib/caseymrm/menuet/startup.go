// +build darwin

package menuet

import (
	"log"
	"os"
	"os/user"
	"sync"
	"text/template"
)

func (a *Application) getStartupPath() string {
	if a.Label == "" {
		log.Fatal("Need to set a Label for the app")
	}
	u, err := user.Current()
	if err != nil {
		log.Printf("user.Current: %v", err)
		return ""
	}
	return u.HomeDir + "/Library/LaunchAgents/" + a.Label + ".plist"
}

func (a *Application) runningAtStartup() bool {
	if a.Label == "" {
		log.Println("Warning: no application Label set")
		return false
	}
	_, err := os.Stat(a.getStartupPath())
	if err == nil {
		return true
	}
	return false
}

func (a *Application) removeStartupItem() {
	err := os.Remove(a.getStartupPath())
	if err != nil {
		log.Printf("os.Remove: %v", err)
	}
}

var launchdOnce sync.Once
var launchdTemplate *template.Template

func (a *Application) addStartupItem() {
	//path := a.getStartupPath()
	//// Make sure ~/Library/LaunchAgents exists
	//err := os.MkdirAll(filepath.Dir(path), 0700)
	//if err != nil {
	//	log.Printf("os.MkdirAll: %v", err)
	//	return
	//}
	//executable, err := os.Executable()
	//if err != nil {
	//	log.Printf("os.Executable: %v", err)
	//	return
	//}
	//f, err := os.Create(path)
	//if err != nil {
	//	log.Printf("os.Create: %v", err)
	//	return
	//}
	//defer f.Close()
	//launchdOnce.Do(func() {
	//	launchdTemplate = template.Must(template.New("launchdConfig").Parse(launchdString))
	//})
	//err = launchdTemplate.Execute(f,
	//	struct {
	//		Name       string
	//		Label      string
	//		Executable string
	//	}{
	//		a.Name,
	//		a.Label,
	//		executable,
	//	})
	//if err != nil {
	//	log.Printf("template.Execute: %v", err)
	//	return
	//}
}

var launchdString = `
<?xml version='1.0' encoding='UTF-8'?>
 <!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
 <plist version='1.0'>
   <dict>
     <key>Label</key><string>{{.Name}}</string>
     <key>Program</key><string>{{.Executable}}</string>
     <key>StandardOutPath</key><string>/tmp/{{.Label}}-out.log</string>
     <key>StandardErrorPath</key><string>/tmp/{{.Label}}-err.log</string>
     <key>RunAtLoad</key><true/>
   </dict>
</plist>
`
