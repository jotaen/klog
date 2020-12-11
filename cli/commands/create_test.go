package commands

import (
	"github.com/stretchr/testify/assert"
	"klog/cli"
	. "klog/testutil/withenv"
	"testing"
)

func TestCreateProject(t *testing.T) {
	WithEnvironment(func(env cli.Environment) {
		code := Create(env, []string{"create"})
		assert.Equal(t, cli.OK, code)
	})
}

func TestCreateProjectAtDate(t *testing.T) {
	WithEnvironment(func(env cli.Environment) {
		code := Create(env, []string{"create", "-d", "1995-03-25"})
		assert.Equal(t, cli.OK, code)
	})
}
