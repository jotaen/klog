package command

import "strings"

type Command struct {
	Bin  string
	Args []string
}

func New(bin string, args []string) Command {
	return Command{Bin: bin, Args: args}
}

func (c Command) ToString() string {
	return c.Bin + " " + strings.Join(c.Args, " ")
}
