package app

import (
	"errors"
	"strings"

	"github.com/jotaen/genie"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/service"
	tf "github.com/jotaen/klog/lib/terminalformat"
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

	// NoWarnings indicates klog should suppress any warning types.
	NoWarnings OptionalParam[service.DisabledCheckers]

	originalConfigFile genie.Data
}

// NewConfig creates a new application configuration by merging the config
// as provided by the various sources.
func NewConfig(numCpus int, getEnvVar func(string) string, configFileContents string) (Config, Error) {
	config := NewDefaultConfig(tf.COLOUR_THEME_DARK)

	config.CpuKernels.set(numCpus)

	// If the config file specifies an `editor` entry as well, the config file
	// takes precedence. This is to allow users to provide a klog-specific
	// editor setting, even if they generally have a `$EDITOR` variable set up.
	if getEnvVar("EDITOR") != "" {
		config.Editor.set(getEnvVar("EDITOR"))
	}

	data, lErr := genie.Parse(configFileContents)
	if lErr != nil {
		return Config{}, NewError(
			"Invalid config file",
			lErr.Error(),
			nil,
		)
	}
	config.originalConfigFile = data
	for _, entry := range CONFIG_FILE_ENTRIES {
		key := entry.Name
		value := data.Get(key)
		if value == "" {
			continue
		}
		rErr := entry.read(value, &config)
		if rErr != nil {
			return Config{}, NewError(
				"Invalid config file",
				"The value for the `"+key+"` setting is not valid: "+entry.Help.Value,
				rErr,
			)
		}
	}

	// If `$NO_COLOR` is set, it takes precedence over the `colour_scheme`
	// entry from the config file. This is to allow users to disable output
	// colouring on the fly, e.g. for programmatic invocation in scripted
	// contexts.
	if getEnvVar("NO_COLOR") != "" {
		config.ColourScheme.set(tf.COLOUR_THEME_NO_COLOUR)
	}
	if getEnvVar("KLOG_DEBUG") != "" {
		config.IsDebug.set(true)
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
		NoWarnings:         newOptionalParam[service.DisabledCheckers](),
	}
}

type Help struct {
	Summary string
	Default string
	Value   string
}

type ConfigFileEntries[T any] struct {
	Name string
	Help Help
	read func(string, *Config) error
}

func (e ConfigFileEntries[T]) Value(c Config) string {
	return c.originalConfigFile.Get(e.Name)
}

var CONFIG_FILE_ENTRIES = []ConfigFileEntries[any]{
	{
		Name: "editor",
		Help: Help{
			Summary: "The CLI command that shall be invoked when running `klog edit`.",
			Value:   "The config property can be any valid CLI command, as you would type it on the terminal. klog will append the target file path as last input argument to that command. Note: you can use quotes in order to prevent undesired shell word-splitting, e.g. if the command name/path contains spaces.",
			Default: "If absent/empty, `klog edit` tries to fall back to the $EDITOR environment variable.",
		},
		read: func(value string, config *Config) error {
			if value != "" {
				config.Editor.set(value)
			}
			return nil
		},
	}, {
		Name: "colour_scheme",
		Help: Help{
			Summary: "The colour scheme of your terminal, so that klog can choose an optimal colour theme for its output.",
			Value:   "The config property must be one of: `dark`, `light`, `basic`, or `no_colour`",
			Default: "If absent/empty, klog assumes a `dark` theme.",
		},
		read: func(value string, config *Config) error {
			switch value {
			case string(tf.COLOUR_THEME_DARK):
				config.ColourScheme.set(tf.COLOUR_THEME_DARK)
			case string(tf.COLOUR_THEME_NO_COLOUR):
				config.ColourScheme.set(tf.COLOUR_THEME_NO_COLOUR)
			case string(tf.COLOUR_THEME_LIGHT):
				config.ColourScheme.set(tf.COLOUR_THEME_LIGHT)
			case string(tf.COLOUR_THEME_BASIC):
				config.ColourScheme.set(tf.COLOUR_THEME_BASIC)
			default:
				return errors.New("The value must be `dark`, `light`, `basic`, or `no_colour`")
			}
			return nil
		},
	}, {
		Name: "default_rounding",
		Help: Help{
			Summary: "The default value that shall be used for rounding input times via the `--round` flag, e.g. in `klog start --round 15m`.",
			Value:   "The config property must be one of: `5m`, `10m`, `15m`, `30m`, `60m`.",
			Default: "If absent/empty, klog doesn’t round input times.",
		},
		read: func(value string, config *Config) error {
			rounding, err := service.NewRoundingFromString(value)
			if err != nil {
				return err
			}
			config.DefaultRounding.set(rounding)
			return nil
		},
	}, {
		Name: "default_should_total",
		Help: Help{
			Summary: "The default duration value that shall be used as should-total when creating new records, e.g. in `klog create --should '8h!'`.",
			Value:   "The config property must be a duration followed by an exclamation mark. Examples: `8h!`, `6h30m!`.",
			Default: "If absent/empty, klog doesn’t set a should-total on new records.",
		},
		read: func(value string, config *Config) error {
			value = strings.TrimSuffix(value, "!")
			d, err := klog.NewDurationFromString(value)
			if err != nil {
				return err
			}
			config.DefaultShouldTotal.set(klog.NewShouldTotal(0, d.InMinutes()))
			return nil
		},
	}, {
		Name: "date_format",
		Help: Help{
			Summary: "The preferred date notation for klog to use when adding a new record to a target file, i.e. whether it uses dashes (as in `2022-03-24`) or slashes (as in `2022/03/24`) as delimiter.",
			Value:   "The config property must be either `YYYY-MM-DD` or `YYYY/MM/DD`.",
			Default: "If absent/empty, klog automatically tries to be consistent with what is used in the target file; in doubt, it defaults to the YYYY-MM-DD format.",
		},
		read: func(value string, config *Config) error {
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
	}, {
		Name: "time_convention",
		Help: Help{
			Summary: "The preferred time convention for klog to use when adding a new time range entry to a target file, i.e. whether it uses the 24-hour clock (as in `13:00`) or the 12-hour clock (as in `1:00pm`).",
			Value:   "The config property must be either `24h` or `12h`.",
			Default: "If absent/empty, klog automatically tries to be consistent with what is used in the target file; in doubt, it defaults to the 24-hour clock format.",
		},
		read: func(value string, config *Config) error {
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
	}, {
		Name: "no_warnings",
		Help: Help{
			Summary: "Whether klog should suppress certain warnings when processing files.",
			Value:   "The config property must be one (or several comma-separated) of: `UNCLOSED_OPEN_RANGE` (for unclosed open ranges in past records), `FUTURE_ENTRIES` (for records/entries in the future), `OVERLAPPING_RANGES` (for time ranges that overlap), `MORE_THAN_24H` (if there is a record with more than 24h total). Multiple values must be separated by a comma, e.g.: `UNCLOSED_OPEN_RANGE, MORE_THAN_24H`.",
			Default: "If absent/empty, klog prints all available warnings.",
		},
		read: func(value string, config *Config) error {
			sanitizedString := strings.ReplaceAll(value, " ", "")
			warningConfigs := strings.Split(sanitizedString, ",")
			disabledCheckers := service.NewDisabledCheckers()
			for _, c := range warningConfigs {
				if _, nameExists := disabledCheckers[c]; !nameExists {
					return errors.New(
						"The value must be a valid warning name, such as `UNCLOSED_OPEN_RANGE`, got: " + c + ".",
					)
				}
				disabledCheckers[c] = true
			}

			config.NoWarnings.set(disabledCheckers)
			return nil
		},
	},
}

type baseParam[T any] struct {
	value T
	isSet bool
}

func (p *baseParam[T]) set(value T) {
	p.value = value
	p.isSet = true
}

type MandatoryParam[T any] struct {
	baseParam[T]
}

func newMandatoryParam[T any](defaultValue T) MandatoryParam[T] {
	return MandatoryParam[T]{baseParam[T]{
		value: defaultValue,
	}}
}

func (p MandatoryParam[T]) Value() T {
	return p.value
}

type OptionalParam[T any] struct {
	baseParam[T]
}

func newOptionalParam[T any]() OptionalParam[T] {
	return OptionalParam[T]{baseParam[T]{
		isSet: false,
	}}
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
