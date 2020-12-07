package cli

import (
	"klog/cli/commands"
	"klog/store"
)

type cmd func(store.Store)

var cmdDict map[string]cmd = map[string]cmd{
	"list":   commands.List,
	"create": commands.Create,
}

func Exec(st store.Store, name string) {
	c := cmdDict[name]
	if c == nil {
		return
	}
	c(st)
}
