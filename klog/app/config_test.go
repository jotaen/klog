package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func createMockConfigFromEnv(vs map[string]string) FromEnvVars {
	return FromEnvVars{GetVar: func(n string) string {
		return vs[n]
	}}
}

func TestCreatesNewDefaultConfig(t *testing.T) {
	c := NewDefaultConfig()
	assert.Equal(t, c.IsDebug.Value(), false)
	assert.Equal(t, c.Editor.Value(), "")
	assert.Equal(t, c.NoColour.Value(), false)
	assert.Equal(t, c.CpuKernels.Value(), 1)
	assert.Equal(t, c.DefaultRounding.Value().ToInt(), 15)
}

func TestSetsParamsMetadataIsHandledCorrectly(t *testing.T) {
	{
		c := NewDefaultConfig()
		assert.Equal(t, c.NoColour.Value(), false)
		assert.False(t, c.NoColour.WasExplicitlySet())
	}
	{
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{
				"NO_COLOR": "1",
			}),
			FromConfigFile{""},
		)
		assert.Equal(t, c.NoColour.Value(), true)
		assert.True(t, c.NoColour.WasExplicitlySet())
	}
}

func TestSetsParamsFromEnv(t *testing.T) {
	// Read plain environment variables.
	{
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{
				"EDITOR":     "subl",
				"KLOG_DEBUG": "1",
				"NO_COLOR":   "1",
			}),
			FromConfigFile{""},
		)
		assert.Equal(t, c.IsDebug.Value(), true)
		assert.Equal(t, c.NoColour.Value(), true)
		assert.Equal(t, c.Editor.Value(), "subl")
	}

	// `KLOG_EDITOR` would trump `EDITOR`.
	{
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{
				"EDITOR":      "subl",
				"KLOG_EDITOR": "vi",
			}),
			FromConfigFile{""},
		)
		assert.Equal(t, c.Editor.Value(), "vi")
	}
}

func TestSetsDefaultRoundingParamFromConfigFile(t *testing.T) {
	for _, tml := range []string{
		`default_rounding: 30m`,
	} {
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.Equal(t, c.DefaultRounding.Value().ToInt(), 30)
		assert.True(t, c.DefaultRounding.WasExplicitlySet())
	}
}

func TestSetsDefaultShouldTotalParamFromConfigFile(t *testing.T) {
	for _, tml := range []string{
		`default_should_total: 8h30m!`,
	} {
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.Equal(t, c.DefaultShouldTotal.Value().ToString(), "8h30m!")
		assert.True(t, c.DefaultShouldTotal.WasExplicitlySet())
	}
}

func TestIgnoresUnknownPropertiesInConfigFile(t *testing.T) {
	for _, tml := range []string{`
unknown_property: 1
what_is_this:
  - 1
  - 2
`,
	} {
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.False(t, c.DefaultRounding.WasExplicitlySet())
	}
}

func TestRejectsInvalidConfigFile(t *testing.T) {
	for _, tml := range []string{
		`default_rounding: true`,              // Wrong type
		`default_rounding: 22m`,               // Invalid value
		`default_should_total: [true, false]`, // Wrong type
		`default_should_total: 15`,            // Invalid value
	} {
		_, err := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.Error(t, err)
	}
}
