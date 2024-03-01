package app

import (
	"errors"
	"github.com/jotaen/genie"
	"github.com/jotaen/klog/klog"
	tf "github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/klog/service"
	"strings"
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
}

type Reader interface {
	Apply(*Config) Error
}

func NewConfig(c1 FromStaticValues, c2 FromEnvVars, c3 FromConfigFile) (Config, Error) {
	config := NewDefaultConfig(tf.DARK)
	for _, c := range []Reader{c1, c2, c3} {
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

type MandatoryParam[T any] struct {
	value        T
	wasSetInFile bool
}

func newMandatoryParam[T any](defaultValue T) MandatoryParam[T] {
	return MandatoryParam[T]{
		value:        defaultValue,
		wasSetInFile: false,
	}
}

func (p MandatoryParam[T]) Value() T {
	return p.value
}

func (p *MandatoryParam[T]) override(value T) {
	p.value = value
}

type OptionalParam[T any] struct {
	value T
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

func (p *OptionalParam[T]) set(value T) {
	p.value = value
	p.isSet = true
}

// FromStaticValues is the part of the configuration that is automatically
// determined, e.g. by constraints of the runtime environment.
type FromStaticValues struct {
	NumCpus int
}

func (e FromStaticValues) Apply(config *Config) Error {
	config.CpuKernels.override(e.NumCpus)
	return nil
}

// FromEnvVars is the part of the configuration that is determined
// by environment variables.
type FromEnvVars struct {
	GetVar func(string) string
}

func (e FromEnvVars) Apply(config *Config) Error {
	if e.GetVar("KLOG_DEBUG") != "" {
		config.IsDebug.override(true)
	}
	if e.GetVar("NO_COLOR") != "" {
		config.ColourScheme.override(tf.NO_COLOUR)
	}
	if e.GetVar("EDITOR") != "" {
		config.Editor.set(e.GetVar("EDITOR"))
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
		Reader: func(value string, config *Config) error {
			if value != "" {
				config.Editor.set(value)
			}
			return nil
		},
		Value: func(c Config) string {
			return c.Editor.value
		},
		Help: Help{
			Summary: "The CLI command that shall be invoked when running `klog edit`.",
			Value:   "The config property can be any valid CLI command, as you would type it on the terminal. klog will append the target file path as last input argument to that command. Note: you can use quotes in order to prevent undesired shell word-splitting, e.g. if the command name/path contains spaces.",
			Default: "If absent/empty, `klog edit` tries to fall back to the $EDITOR environment variable (which you’d see below in that case).",
		},
	}, {
		Name: "colour_scheme",
		Reader: func(value string, config *Config) error {
			switch value {
			case string(tf.DARK):
				config.ColourScheme.override(tf.DARK)
				config.ColourScheme.wasSetInFile = true
			case string(tf.LIGHT):
				config.ColourScheme.override(tf.LIGHT)
				config.ColourScheme.wasSetInFile = true
			case string(tf.NO_COLOUR):
				config.ColourScheme.override(tf.NO_COLOUR)
				config.ColourScheme.wasSetInFile = true
			default:
				return errors.New("The value must be `dark`, `light` or `no_colour`")
			}
			return nil
		},
		Value: func(c Config) string {
			if !c.ColourScheme.wasSetInFile {
				return ""
			}
			return string(c.ColourScheme.Value())
		},
		Help: Help{
			Summary: "The colour scheme of your terminal, so that klog can choose an optimal colour theme for its output.",
			Value:   "The config property must be one of: `dark`, `light` or `no_colour`",
			Default: "If absent/empty, klog assumes a `dark` theme.",
		},
	}, {
		Name: "default_rounding",
		Reader: func(value string, config *Config) error {
			rounding, err := service.NewRoundingFromString(value)
			if err != nil {
				return err
			}
			config.DefaultRounding.set(rounding)
			return nil
		},
		Value: func(c Config) string {
			result := ""
			c.DefaultRounding.Unwrap(func(r service.Rounding) {
				result = r.ToString()
			})
			return result
		},
		Help: Help{
			Summary: "The default value that shall be used for rounding input times via the `--round` flag, e.g. in `klog start --round 15m`.",
			Value:   "The config property must be one of: `5m`, `10m`, `15m`, `30m`, `60m`.",
			Default: "If absent/empty, klog doesn’t round input times.",
		},
	}, {
		Name: "default_should_total",
		Reader: func(value string, config *Config) error {
			value = strings.TrimSuffix(value, "!")
			d, err := klog.NewDurationFromString(value)
			if err != nil {
				return err
			}
			config.DefaultShouldTotal.set(klog.NewShouldTotal(0, d.InMinutes()))
			return nil
		},
		Value: func(c Config) string {
			result := ""
			c.DefaultShouldTotal.Unwrap(func(s klog.ShouldTotal) {
				result = s.ToString()
			})
			return result
		},
		Help: Help{
			Summary: "The default duration value that shall be used as should-total when creating new records, e.g. in `klog create --should '8h!'`.",
			Value:   "The config property must be a duration followed by an exclamation mark. Examples: `8h!`, `6h30m!`.",
			Default: "If absent/empty, klog doesn’t set a should-total on new records.",
		},
	}, {
		Name: "date_format",
		Reader: func(value string, config *Config) error {
			useDashes := true
			if value == "YYYY-MM-DD" {
				useDashes = true
			} else if value == "YYYY/MM/DD" {
				useDashes = false
			} else {
				return errors.New("The value must be `YYYY-MM-DD` or `YYYY/MM/DD`")
			}
			config.DateUseDashes.set(useDashes)
			return nil
		},
		Value: func(c Config) string {
			result := ""
			c.DateUseDashes.Unwrap(func(d bool) {
				if d {
					result = "YYYY-MM-DD"
				} else {
					result = "YYYY/MM/DD"
				}
			})
			return result
		},
		Help: Help{
			Summary: "The preferred date notation for klog to use when adding a new record to a target file, i.e. whether it uses dashes (as in `2022-03-24`) or slashes (as in `2022/03/24`) as delimiter.",
			Value:   "The config property must be either `YYYY-MM-DD` or `YYYY/MM/DD`.",
			Default: "If absent/empty, klog automatically tries to be consistent with what is used in the target file; in doubt, it defaults to the YYYY-MM-DD format.",
		},
	}, {
		Name: "time_convention",
		Reader: func(value string, config *Config) error {
			use24HourClock := true
			if value == "24h" {
				use24HourClock = true
			} else if value == "12h" {
				use24HourClock = false
			} else {
				return errors.New("The value must be `24h` or `12h`")
			}
			config.TimeUse24HourClock.set(use24HourClock)
			return nil
		},
		Value: func(c Config) string {
			result := ""
			c.TimeUse24HourClock.Unwrap(func(t bool) {
				if t {
					result = "24h"
				} else {
					result = "12h"
				}
			})
			return result
		},
		Help: Help{
			Summary: "The preferred time convention for klog to use when adding a new time range entry to a target file, i.e. whether it uses the 24-hour clock (as in `13:00`) or the 12-hour clock (as in `1:00pm`).",
			Value:   "The config property must be either `24h` or `12h`.",
			Default: "If absent/empty, klog automatically tries to be consistent with what is used in the target file; in doubt, it defaults to the 24-hour clock format.",
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
	Reader func(string, *Config) error
	Value  func(Config) string
	Help   Help
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
		rErr := entry.Reader(value, config)
		if rErr != nil {
			return NewError(
				"Invalid config file",
				"The value for `"+key+"` is not valid: "+entry.Help.Value,
				rErr,
			)
		}

	}
	return nil
}
