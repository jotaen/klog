package cli

import (
	"klog/app"
	systray "klog/app/mac_widget"
)

type Widget struct {
}

func (args *Print) Widget(ctx *app.Context) error {
	systray.Launch()
	return nil
}
