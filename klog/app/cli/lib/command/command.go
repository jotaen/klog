package command

import (
	"errors"
	"github.com/kballard/go-shellquote"
)

type Command struct {
	Bin  string
	Args []string
}

func NewFromString(command string) (Command, error) {
	words, err := shellquote.Split(command)
	if err != nil {
		return Command{}, err
	}
	if len(words) == 0 {
		return Command{}, errors.New("Empty command")
	}
	return New(words[0], words[1:]), nil
}

func New(bin string, args []string) Command {
	return Command{Bin: bin, Args: args}
}
