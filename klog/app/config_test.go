package app

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/service"
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

	isRoundingSet := false
	c.DefaultRounding.Map(func(_ service.Rounding) {
		isRoundingSet = true
	})
	assert.False(t, isRoundingSet)

	isShouldTotalSet := false
	c.DefaultShouldTotal.Map(func(_ klog.ShouldTotal) {
		isShouldTotalSet = true
	})
	assert.False(t, isShouldTotalSet)
}

func TestSetsParamsMetadataIsHandledCorrectly(t *testing.T) {
	{
		c := NewDefaultConfig()
		assert.Equal(t, c.NoColour.Value(), false)
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
		`default_rounding = 30m`,
	} {
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		var value int
		c.DefaultRounding.Map(func(r service.Rounding) {
			value = r.ToInt()
		})
		assert.Equal(t, value, 30)
	}
}

func TestSetsDefaultShouldTotalParamFromConfigFile(t *testing.T) {
	for _, tml := range []string{
		`default_should_total = 8h30m!`,
	} {
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		var value string
		c.DefaultShouldTotal.Map(func(s klog.ShouldTotal) {
			value = s.ToString()
		})
		assert.Equal(t, value, "8h30m!")
	}
}

func TestIgnoresUnknownPropertiesInConfigFile(t *testing.T) {
	for _, tml := range []string{`
unknown_property = 1
what_is_this = true
`,
	} {
		_, err := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.Nil(t, err)
	}
}

func TestIgnoresEmptyConfigFileOrEmptyParameters(t *testing.T) {
	for _, tml := range []string{
		``,
		`
default_rounding =
`,
		`
default_should_total = 
`,
	} {
		_, err := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.Nil(t, err)
	}
}

func TestRejectsInvalidConfigFile(t *testing.T) {
	for _, tml := range []string{
		`default_rounding = true`,              // Wrong type
		`default_rounding = 22m`,               // Invalid value
		`default_should_total = [true, false]`, // Wrong type
		`default_should_total = 15`,            // Invalid value
	} {
		_, err := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.Error(t, err)
	}
}
