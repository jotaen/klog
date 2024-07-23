package app

import (
	"errors"
	"strings"

	"github.com/jotaen/genie"
	"github.com/jotaen/klog/klog"
	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/service"
)

// Config contain all variable settings that influence the behaviour of
// the application. Some of these properties can be controlled by the user
// in one way or the other, depending on their purpose.
type Config struct {
	// IsDebug indicates whether klog should print additional debug information.
	// This is an ephemeral property, which is used for debugging purposes, and not
	// supposed to be configured permanently.
	IsDebug MandatoryParam[bool]

	// Editor is the CLI command with which to invoke the preferred editor.
	Editor OptionalParam[string]

	// ColourScheme specifies the background of the terminal, so that
	// the output colours can be adjusted accordingly.
	ColourScheme MandatoryParam[tf.ColourTheme]

	// CpuKernels is the number of available CPUs that klog is allowed to utilise.
	// The value must be `1` or higher.
	// This is a low-level property that is not supposed to be exposed to end-users at all.
	CpuKernels MandatoryParam[int]

	// DefaultRounding is the default for the --round flag.
	DefaultRounding OptionalParam[service.Rounding]

	// DefaultShouldTotal is the default should total for new records.
	DefaultShouldTotal OptionalParam[klog.ShouldTotal]

	// DateUseDashes denotes the preferred date format: YYYY-MM-DD (true) or YYYY/MM/DD (false).
	DateUseDashes OptionalParam[bool]

	// TimeUse24HourClock denotes the preferred time format: 13:00 (true) or 1:00pm (false).
	TimeUse24HourClock OptionalParam[bool]

	// HideWarnings indicates klog should suppress any warning types.
	HideWarnings OptionalParam[[]string]
}

type Reader interface {
	Apply(*Config) Error
}

// NewConfig creates a new application configuration by merging the config
// based on the following precedence: (1) env variables, (2) config file,
// (3) determined values.
func NewConfig(determined FromDeterminedValues, env FromEnvVars, file FromConfigFile) (Config, Error) {
	config := NewDefaultConfig(tf.COLOUR_THEME_DARK)
	for _, c := range []Reader{determined, file, env} {
		err := c.Apply(&config)
		if err != nil {
			return Config{}, err
		}
	}
	return config, nil
}

func NewDefaultConfig(c tf.ColourTheme) Config {
	return Config{
		IsDebug:            newMandatoryParam(false),
		Editor:             newOptionalParam[string](),
		ColourScheme:       newMandatoryParam(c),
		CpuKernels:         newMandatoryParam(1),
		DefaultRounding:    newOptionalParam[service.Rounding](),
		DefaultShouldTotal: newOptionalParam[klog.ShouldTotal](),
	}
}

type configOrigin int

const (
	configOriginEnv = iota + 1
	configOriginFile
	configOriginStaticValues
)

type BaseParam[T any] struct {
	value  T
	origin configOrigin
}

type MandatoryParam[T any] struct {
	BaseParam[T]
}

func newMandatoryParam[T any](defaultValue T) MandatoryParam[T] {
	return MandatoryParam[T]{BaseParam[T]{
		value:  defaultValue,
		origin: 0,
	}}
}

func (p MandatoryParam[T]) Value() T {
	return p.value
}

func (p *MandatoryParam[T]) override(value T, o configOrigin) {
	p.value = value
	p.origin = o
}

type OptionalParam[T any] struct {
	BaseParam[T]
	isSet bool
}

func newOptionalParam[T any]() OptionalParam[T] {
	return OptionalParam[T]{
		isSet: false,
	}
}

func (p OptionalParam[T]) Unwrap(f func(T)) {
	if p.isSet {
		f(p.value)
	}
}

func (p OptionalParam[T]) UnwrapOr(defaultValue T) T {
	if p.isSet {
		return p.value
	}
	return defaultValue
}

func (p *OptionalParam[T]) set(value T, o configOrigin) {
	p.value = value
	p.isSet = true
	p.origin = o
}

// FromDeterminedValues is the part of the configuration that is automatically
// determined, e.g. by constraints of the runtime environment.
type FromDeterminedValues struct {
	NumCpus int
}

func (e FromDeterminedValues) Apply(config *Config) Error {
	config.CpuKernels.override(e.NumCpus, configOriginStaticValues)
	return nil
}

// FromEnvVars is the part of the configuration that is determined
// by environment variables.
type FromEnvVars struct {
	GetVar func(string) string
}

func (e FromEnvVars) Apply(config *Config) Error {
	if e.GetVar("KLOG_DEBUG") != "" {
		config.IsDebug.override(true, configOriginEnv)
	}
	if e.GetVar("NO_COLOR") != "" {
		config.ColourScheme.override(tf.COLOUR_THEME_NO_COLOUR, configOriginEnv)
	}
	if e.GetVar("EDITOR") != "" {
		config.Editor.set(e.GetVar("EDITOR"), configOriginEnv)
	}
	return nil
}

// FromConfigFile is the part of the configuration that the user can
// control via a configuration file.
type FromConfigFile struct {
	FileContents string
}

var CONFIG_FILE_ENTRIES = []ConfigFileEntries[any]{
	{
		Name: "editor",
		reader: func(value string, config *Config) error {
			if value != "" {
				config.Editor.set(value, configOriginFile)
			}
			return nil
		},
		value: func(c Config) (string, configOrigin) {
			return c.Editor.value, c.Editor.origin
		},
		Help: Help{
			Summary: "The CLI command that shall be invoked when running `klog edit`.",
			Value:   "The config property can be any valid CLI command, as you would type it on the terminal. klog will append the target file path as last input argument to that command. Note: you can use quotes in order to prevent undesired shell word-splitting, e.g. if the command name/path contains spaces.",
			Default: "If absent/empty, `klog edit` tries to fall back to the $EDITOR environment variable (which, by the way, takes precedence, if set).",
		},
	}, {
		Name: "colour_scheme",
		reader: func(value string, config *Config) error {
			switch value {
			case string(tf.COLOUR_THEME_DARK):
				config.ColourScheme.override(tf.COLOUR_THEME_DARK, configOriginFile)
			case string(tf.COLOUR_THEME_NO_COLOUR):
				config.ColourScheme.override(tf.COLOUR_THEME_NO_COLOUR, configOriginFile)
			case string(tf.COLOUR_THEME_LIGHT):
				config.ColourScheme.override(tf.COLOUR_THEME_LIGHT, configOriginFile)
			default:
				return errors.New("The value must be `dark`, `light` or `no_colour`")
			}
			return nil
		},
		value: func(c Config) (string, configOrigin) {
			return string(c.ColourScheme.Value()), c.ColourScheme.origin
		},
		Help: Help{
			Summary: "The colour scheme of your terminal, so that klog can choose an optimal colour theme for its output.",
			Value:   "The config property must be one of: `dark`, `light` or `no_colour`",
			Default: "If absent/empty, klog assumes a `dark` theme.",
		},
	}, {
		Name: "default_rounding",
		reader: func(value string, config *Config) error {
			rounding, err := service.NewRoundingFromString(value)
			if err != nil {
				return err
			}
			config.DefaultRounding.set(rounding, configOriginFile)
			return nil
		},
		value: func(c Config) (string, configOrigin) {
			result := ""
			c.DefaultRounding.Unwrap(func(r service.Rounding) {
				result = r.ToString()
			})
			return result, c.DefaultRounding.origin
		},
		Help: Help{
			Summary: "The default value that shall be used for rounding input times via the `--round` flag, e.g. in `klog start --round 15m`.",
			Value:   "The config property must be one of: `5m`, `10m`, `15m`, `30m`, `60m`.",
			Default: "If absent/empty, klog doesn’t round input times.",
		},
	}, {
		Name: "default_should_total",
		reader: func(value string, config *Config) error {
			value = strings.TrimSuffix(value, "!")
			d, err := klog.NewDurationFromString(value)
			if err != nil {
				return err
			}
			config.DefaultShouldTotal.set(klog.NewShouldTotal(0, d.InMinutes()), configOriginFile)
			return nil
		},
		value: func(c Config) (string, configOrigin) {
			result := ""
			c.DefaultShouldTotal.Unwrap(func(s klog.ShouldTotal) {
				result = s.ToString()
			})
			return result, c.DefaultShouldTotal.origin
		},
		Help: Help{
			Summary: "The default duration value that shall be used as should-total when creating new records, e.g. in `klog create --should '8h!'`.",
			Value:   "The config property must be a duration followed by an exclamation mark. Examples: `8h!`, `6h30m!`.",
			Default: "If absent/empty, klog doesn’t set a should-total on new records.",
		},
	}, {
		Name: "date_format",
		reader: func(value string, config *Config) error {
			useDashes := true
			if value == "YYYY-MM-DD" {
				useDashes = true
			} else if value == "YYYY/MM/DD" {
				useDashes = false
			} else {
				return errors.New("The value must be `YYYY-MM-DD` or `YYYY/MM/DD`")
			}
			config.DateUseDashes.set(useDashes, configOriginFile)
			return nil
		},
		value: func(c Config) (string, configOrigin) {
			result := ""
			c.DateUseDashes.Unwrap(func(d bool) {
				if d {
					result = "YYYY-MM-DD"
				} else {
					result = "YYYY/MM/DD"
				}
			})
			return result, c.DateUseDashes.origin
		},
		Help: Help{
			Summary: "The preferred date notation for klog to use when adding a new record to a target file, i.e. whether it uses dashes (as in `2022-03-24`) or slashes (as in `2022/03/24`) as delimiter.",
			Value:   "The config property must be either `YYYY-MM-DD` or `YYYY/MM/DD`.",
			Default: "If absent/empty, klog automatically tries to be consistent with what is used in the target file; in doubt, it defaults to the YYYY-MM-DD format.",
		},
	}, {
		Name: "time_convention",
		reader: func(value string, config *Config) error {
			use24HourClock := true
			if value == "24h" {
				use24HourClock = true
			} else if value == "12h" {
				use24HourClock = false
			} else {
				return errors.New("The value must be `24h` or `12h`")
			}
			config.TimeUse24HourClock.set(use24HourClock, configOriginFile)
			return nil
		},
		value: func(c Config) (string, configOrigin) {
			result := ""
			c.TimeUse24HourClock.Unwrap(func(t bool) {
				if t {
					result = "24h"
				} else {
					result = "12h"
				}
			})
			return result, c.TimeUse24HourClock.origin
		},
		Help: Help{
			Summary: "The preferred time convention for klog to use when adding a new time range entry to a target file, i.e. whether it uses the 24-hour clock (as in `13:00`) or the 12-hour clock (as in `1:00pm`).",
			Value:   "The config property must be either `24h` or `12h`.",
			Default: "If absent/empty, klog automatically tries to be consistent with what is used in the target file; in doubt, it defaults to the 24-hour clock format.",
		},
	}, {
		Name: "no_warnings",
		reader: func(value string, config *Config) error {
			sanitizedString := strings.ToLower(strings.ReplaceAll(value, " ", ""))
			warningConfigs := strings.Split(sanitizedString, ",")

			hideWarnings := []string{}
			for _, warningConfig := range warningConfigs {
				checkerName, err := convertToCheckerName(warningConfig)
				if err != nil {
					return err
				}
				hideWarnings = append(hideWarnings, checkerName)
			}

			config.HideWarnings.set(hideWarnings, configOriginFile)
			return nil
		},
		value: func(c Config) (string, configOrigin) {
			result := ""
			c.HideWarnings.Unwrap(func(warningConfigs []string) {
				result = strings.Join(warningConfigs, " , ")
			})
			return result, c.HideWarnings.origin
		},
		Help: Help{
			Summary: "Whether klog should suppress warnings when printing time reports.",
			Value:   "This value can be left blank, or set to one or more of the following: UNLCLOSED_RANGE, FUTURE_ENTRY, OVERLAPPED_RANGE, GREATER_THAN_24HRS. Multiple values must be separated by a comma, e.g. `UNLCLOSED_RANGE,OVERLAPPED_RANGE`.",
			Default: "If absent/empty, klog prints all available warnings.",
		},
	},
}

type Help struct {
	Summary string
	Default string
	Value   string
}

type ConfigFileEntries[T any] struct {
	Name   string
	Help   Help
	reader func(string, *Config) error
	value  func(Config) (string, configOrigin)
}

func (e ConfigFileEntries[T]) Value(c Config) string {
	v, o := e.value(c)
	if o == configOriginFile {
		return v
	}
	return ""
}

func (e FromConfigFile) Apply(config *Config) Error {
	data, lErr := genie.Parse(e.FileContents)
	if lErr != nil {
		return NewError(
			"Invalid config file",
			lErr.Error(),
			nil,
		)
	}
	for _, entry := range CONFIG_FILE_ENTRIES {
		key := entry.Name
		value := data.Get(key)
		if value == "" {
			continue
		}
		rErr := entry.reader(value, config)
		if rErr != nil {
			return NewError(
				"Invalid config file",
				"The value for the `"+key+"` setting is not valid: "+entry.Help.Value,
				rErr,
			)
		}

	}
	return nil
}

// We'd like to disconnect our internal names from those that are exposed to the
// user via the configuration. This function allows us to capture that mapping as
// needed. See the no_warnings entry in CONFIG_FILE_ENTRIES.
func convertToCheckerName(configString string) (string, error) {
	configString = strings.TrimSpace(configString)
	configString = strings.ToLower(configString)
	switch configString {
	case "unclosed_range":
		return "unclosedOpenRange", nil
	case "future_entry":
		return "futureEntries", nil
	case "overlapped_range":
		return "overlappingTimeRanges", nil
	case "greater_than_24hrs":
		return "moreThan24Hours", nil
	default:
		return "", errors.New("Unknown checker name:" + configString)
	}
}
