// +build !darwin

package mac_widget

import "fmt"

func Run(forceRunThroughLaunchAgent bool) {
	fmt.Println("The widget is currently only supported on MacOS.")
}
