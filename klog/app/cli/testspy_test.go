package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/command"
)

type commandSpy struct {
	LastCmd command.Command
	Count   int
	onCmd   func(command.Command) app.Error
}

func (c *commandSpy) Execute(cmd command.Command) app.Error {
	c.LastCmd = cmd
	c.Count++
	return c.onCmd(cmd)
}

func newCommandSpy(onCmd func(command.Command) app.Error) *commandSpy {
	if onCmd == nil {
		onCmd = func(_ command.Command) app.Error {
			return nil
		}
	}
	return &commandSpy{
		LastCmd: command.Command{},
		Count:   0,
		onCmd:   onCmd,
	}
}
