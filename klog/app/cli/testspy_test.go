package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/lib/shellcmd"
)

type commandSpy struct {
	LastCmd shellcmd.Command
	Count   int
	onCmd   func(shellcmd.Command) app.Error
}

func (c *commandSpy) Execute(cmd shellcmd.Command) app.Error {
	c.LastCmd = cmd
	c.Count++
	return c.onCmd(cmd)
}

func newCommandSpy(onCmd func(shellcmd.Command) app.Error) *commandSpy {
	if onCmd == nil {
		onCmd = func(_ shellcmd.Command) app.Error {
			return nil
		}
	}
	return &commandSpy{
		LastCmd: shellcmd.Command{},
		Count:   0,
		onCmd:   onCmd,
	}
}
