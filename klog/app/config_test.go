package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func createMockConfigFromEnv(vs map[string]string) ConfigFromEnvVars {
	return ConfigFromEnvVars{GetVar: func(n string) string {
		return vs[n]
	}}
}

func TestCreatesNewDefaultConfig(t *testing.T) {
	c := NewDefaultConfig()
	assert.Equal(t, c.IsDebug.Value(), false)
	assert.Equal(t, c.Editor.Value(), "")
	assert.Equal(t, c.NoColour.Value(), false)
	assert.Equal(t, c.CpuKernels.Value(), 1)
}

func TestSetsParamsMetadataIsHandledCorrectly(t *testing.T) {
	{
		c := NewDefaultConfig()
		assert.Equal(t, c.NoColour.Value(), false)
		assert.False(t, c.NoColour.WasExplicitlySet())
		assert.Equal(t, c.NoColour.Default(), false)
	}
	{
		c := NewConfig(
			ConfigFromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{
				"NO_COLOR": "1",
			}),
		)
		assert.Equal(t, c.NoColour.Value(), true)
		assert.True(t, c.NoColour.WasExplicitlySet())
		assert.Equal(t, c.NoColour.Default(), false)
	}
}

func TestSetsParamsFromEnv(t *testing.T) {
	// Read plain environment variables.
	{
		c := NewConfig(ConfigFromStaticValues{NumCpus: 1}, createMockConfigFromEnv(map[string]string{
			"EDITOR":     "subl",
			"KLOG_DEBUG": "1",
			"NO_COLOR":   "1",
		}))
		assert.Equal(t, c.IsDebug.Value(), true)
		assert.Equal(t, c.NoColour.Value(), true)
		assert.Equal(t, c.Editor.Value(), "subl")
	}

	// `KLOG_EDITOR` would trump `EDITOR`.
	{
		c := NewConfig(ConfigFromStaticValues{NumCpus: 1}, createMockConfigFromEnv(map[string]string{
			"EDITOR":      "subl",
			"KLOG_EDITOR": "vi",
		}))
		assert.Equal(t, c.Editor.Value(), "vi")
	}
}
