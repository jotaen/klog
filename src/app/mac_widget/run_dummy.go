//go:build !(darwin && amd64)

package mac_widget

func IsWidgetAvailable() bool {
	return false
}

func Run(forceRunThroughLaunchAgent bool) {}
