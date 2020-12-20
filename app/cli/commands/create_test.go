package commands

import (
	"github.com/stretchr/testify/assert"
	"klog/app"
	"klog/app/cli"
	"klog/project"
	. "klog/testutil/withenv"
	"testing"
)

func TestCreateProject(t *testing.T) {
	WithEnvironment(func(env app.Environment, prj project.Project) {
		code := Create.Main(env, prj, []string{"create"})
		assert.Equal(t, cli.OK, code)
	})
}

func TestCreateProjectAtDate(t *testing.T) {
	WithEnvironment(func(env app.Environment, prj project.Project) {
		code := Create.Main(env, prj, []string{"create", "-d", "1995-03-25"})
		assert.Equal(t, cli.OK, code)
	})
}
