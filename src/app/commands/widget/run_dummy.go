//go:build !(darwin && amd64)

package widget

func IsWidgetAvailable() bool {
	return false
}

func Run(forceRunThroughLaunchAgent bool) {}
