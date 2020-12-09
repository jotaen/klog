package cli

import (
	"github.com/stretchr/testify/assert"
	"klog/cli/lib"
	"os"
	"testing"
)

func run(fn func(string)) {
	path := "../tmp/test"
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModePerm)
	fn(path)
	os.RemoveAll(path)
}

func TestErrorIfPathDoesNotExist(t *testing.T) {
	run(func(path string) {
		code := Execute(path + "asdf1234", []string{"create"})
		assert.Equal(t, lib.PROJECT_PATH_INVALID, code)
	})
}

func TestCreateProject(t *testing.T) {
	run(func(path string) {
		code := Execute(path, []string{"create"})
		assert.Equal(t, lib.OK, code)
	})
}
