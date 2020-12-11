package core

import (
	"github.com/stretchr/testify/assert"
	"klog/cli"
	. "klog/testutil/withdisk"
	"testing"
)

func TestErrorIfPathDoesNotExist(t *testing.T) {
	WithDisk(func(path string) {
		code := Execute(path+"asdf1234", []string{"create"})
		assert.Equal(t, cli.PROJECT_PATH_INVALID, code)
	})
}
