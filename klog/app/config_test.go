package app

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
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
	c := NewDefaultConfig(terminalformat.NO_COLOUR)
	assert.Equal(t, c.IsDebug.Value(), false)
	assert.Equal(t, c.Editor.UnwrapOr(""), "")
	assert.Equal(t, c.CpuKernels.Value(), 1)

	isRoundingSet := false
	c.DefaultRounding.Unwrap(func(_ service.Rounding) {
		isRoundingSet = true
	})
	assert.False(t, isRoundingSet)

	isShouldTotalSet := false
	c.DefaultShouldTotal.Unwrap(func(_ klog.ShouldTotal) {
		isShouldTotalSet = true
	})
	assert.False(t, isShouldTotalSet)
}

func TestSetsParamsMetadataIsHandledCorrectly(t *testing.T) {
	{
		c := NewDefaultConfig(terminalformat.NO_COLOUR)
		assert.Equal(t, c.IsDebug.Value(), false)
	}
	{
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{
				"KLOG_DEBUG": "1",
			}),
			FromConfigFile{""},
		)
		assert.Equal(t, c.IsDebug.Value(), true)
	}
}

func TestSetsParamsFromEnv(t *testing.T) {
	t.Run("Read plain environment variables.", func(t *testing.T) {
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
		assert.Equal(t, c.ColourScheme.Value(), terminalformat.NO_COLOUR)
		assert.Equal(t, c.Editor.UnwrapOr(""), "subl")
	})

	t.Run("`editor` from file would trump `$EDITOR` env variable.", func(t *testing.T) {
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{
				"EDITOR": "subl",
			}),
			FromConfigFile{"editor = vi"},
		)
		assert.Equal(t, "vi", c.Editor.UnwrapOr(""))
	})
}

func TestSetsDefaultRoundingParamFromConfigFile(t *testing.T) {
	for _, x := range []struct {
		cfg string
		exp int
	}{
		{`default_rounding = 5m`, 5},
		{`default_rounding = 10m`, 10},
		{`default_rounding = 15m`, 15},
		{`default_rounding = 30m`, 30},
		{`default_rounding = 60m`, 60},
	} {
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{x.cfg},
		)
		var value int
		c.DefaultRounding.Unwrap(func(r service.Rounding) {
			value = r.ToInt()
		})
		assert.Equal(t, x.exp, value)
	}
}

func TestSetsDefaultShouldTotalParamFromConfigFile(t *testing.T) {
	for _, x := range []struct {
		cfg string
		exp string
	}{
		{`default_should_total = 8h30m!`, "8h30m!"},
	} {
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{x.cfg},
		)
		var value string
		c.DefaultShouldTotal.Unwrap(func(s klog.ShouldTotal) {
			value = s.ToString()
		})
		assert.Equal(t, x.exp, value)
	}
}

func TestSetsDateFormatParamFromConfigFile(t *testing.T) {
	for _, x := range []struct {
		cfg string
		exp bool
	}{
		{`date_format = YYYY-MM-DD`, true},
		{`date_format = YYYY/MM/DD`, false},
	} {
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{x.cfg},
		)
		var value bool
		c.DateUseDashes.Unwrap(func(s bool) {
			value = s
		})
		assert.Equal(t, x.exp, value)
	}
}

func TestSetTimeFormatParamFromConfigFile(t *testing.T) {
	for _, x := range []struct {
		cfg string
		exp bool
	}{
		{`time_convention = 24h`, true},
		{`time_convention = 12h`, false},
	} {
		c, _ := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{x.cfg},
		)
		var value bool
		c.TimeUse24HourClock.Unwrap(func(s bool) {
			value = s
		})
		assert.Equal(t, x.exp, value)
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
		`editor = `,
		`default_rounding =`,
		`default_should_total = `,
		`date_format = `,
		`time_convention = `,
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
		`default_rounding = 25m`,               // Invalid value
		`default_should_total = [true, false]`, // Wrong type
		`default_should_total = 15`,            // Invalid value
		`date_format = [true, false]`,          // Wrong type
		`date_format = YYYY.MM.DD`,             // Invalid value
		`time_convention = [true, false]`,      // Wrong type
		`time_convention = 2h`,                 // Invalid value
	} {
		_, err := NewConfig(
			FromStaticValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.Error(t, err)
	}
}
