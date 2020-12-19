package cli

import (
	"klog/project"
)

type Environment struct {
	WorkDir string
	Store   project.Project
}
