package cli

import (
	"klog/app"
)

type Config struct {
	Init bool `name:"sort" help:"Initialise new configuration"`
}

func (args *Config) Run(ctx *app.Context) error {
	return nil
}
