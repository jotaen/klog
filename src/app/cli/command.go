package cli

import (
	"klog/app"
)

type Command struct {
	Main        func(app.Service, []string) int
	Name        string
	Description string
}
