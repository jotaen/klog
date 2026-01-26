package shellcmd

import (
	"errors"

	"github.com/kballard/go-shellquote"
)

// Command represents a shell command invocation with target binary
// and input arguments.
type Command struct {
	Bin  string
	Args []string
}

// NewFromString parses a command invocation into a Command.
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

// New constructs a Command object from values.
func New(bin string, args []string) Command {
	return Command{Bin: bin, Args: args}
}
