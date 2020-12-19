package exec

import (
	"github.com/stretchr/testify/assert"
	"klog/app/cli"
	. "klog/testutil/withdisk"
	"testing"
)

func TestErrorIfPathDoesNotExist(t *testing.T) {
	WithDisk(func(path string) {
		code := Execute(path+"asdf1234", []string{"create"})
		assert.Equal(t, cli.PROJECT_PATH_INVALID, code)
	})
}

func TestErrorSubcommandNotExist(t *testing.T) {
	WithDisk(func(path string) {
		code := Execute(path, []string{"aus6dfri6asydfh"})
		assert.Equal(t, cli.SUBCOMMAND_NOT_FOUND, code)
	})
}
