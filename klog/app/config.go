package app

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/service"
	"gopkg.in/yaml.v3"
	"strings"
)

// Config are all aspects and behaviour of the application that can be
// controlled by the user.
type Config struct {
	// IsDebug indicates whether klog should print additional debug information.
	IsDebug Param[bool]

	// Editor is the CLI command with which to invoke the preferred editor.
	Editor Param[string]

	// NoColour specifies whether output should be coloured.
	NoColour Param[bool]

	// CpuKernels is the number of available CPUs that klog is allowed to utilise.
	// The value must be `1` or higher.
	CpuKernels Param[int]

	// DefaultRounding is the default for the --round flag.
	DefaultRounding Param[service.Rounding]

	// DefaultShouldTotal is the default should total for new records.
	DefaultShouldTotal Param[klog.ShouldTotal]
}

type Reader interface {
	Apply(*Config) Error
}

func NewConfig(c1 FromStaticValues, c2 FromEnvVars, c3 FromConfigFile) (Config, Error) {
	config := NewDefaultConfig()
	for _, c := range []Reader{c1, c2, c3} {
		err := c.Apply(&config)
		if err != nil {
			return Config{}, err
		}
	}
	return config, nil
}

func NewDefaultConfig() Config {
	defaultRounding, err := service.NewRounding(15)
	if err != nil {
		panic(err) // This can/should never happen
	}
	return Config{
		IsDebug:            newDefaultParam(false),
		Editor:             newDefaultParam(""),
		NoColour:           newDefaultParam(false),
		CpuKernels:         newDefaultParam(1),
		DefaultRounding:    newDefaultParam(defaultRounding),
		DefaultShouldTotal: newDefaultParam(klog.NewShouldTotal(0, 0)),
	}
}

type Param[T any] struct {
	actualValue T
	isExplicit  bool
}

func newDefaultParam[T any](value T) Param[T] {
	return Param[T]{
		actualValue: value,
		isExplicit:  false,
	}
}

func (p Param[T]) Value() T {
	return p.actualValue
}

func (p Param[T]) WasExplicitlySet() bool {
	return p.isExplicit
}

func (p *Param[T]) set(value T) {
	p.actualValue = value
	p.isExplicit = true
}

// FromStaticValues is the part of the configuration that is automatically
// determined, e.g. by constraints of the runtime environment.
type FromStaticValues struct {
	NumCpus int
}

func (e FromStaticValues) Apply(config *Config) Error {
	config.CpuKernels.set(e.NumCpus)
	return nil
}

// FromEnvVars is the part of the configuration that is determined
// by environment variables.
type FromEnvVars struct {
	GetVar func(string) string
}

func (e FromEnvVars) Apply(config *Config) Error {
	if e.GetVar("KLOG_DEBUG") != "" {
		config.IsDebug.set(true)
	}
	if e.GetVar("NO_COLOR") != "" {
		config.NoColour.set(true)
	}
	if e.GetVar("KLOG_EDITOR") != "" {
		config.Editor.set(e.GetVar("KLOG_EDITOR"))
	} else if e.GetVar("EDITOR") != "" {
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
			if !c.DefaultRounding.WasExplicitlySet() {
				return ""
			}
			return c.DefaultRounding.Value().ToString()
		},
		Description:  "The default value that shall be used for rounding time values via the `--round` flag, e.g. in `klog start --round 15m`. (If absent/empty, klog doesn’t round.)",
		Instructions: "The value must be one of: `5m`, `10m`, `15m`, `20m`, `30m`, `60m`. ",
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
			if !c.DefaultShouldTotal.WasExplicitlySet() {
				return ""
			}
			return c.DefaultShouldTotal.Value().ToString()
		},
		Description:  "The default duration value that shall be used as should-total when creating new records. (If absent/empty, klog doesn’t set a should-total on new records.)",
		Instructions: "The value must be a duration followed by an exclamation mark. Examples: `8h!`, `6h30m!`. ",
	},
}

type ConfigFileEntries[T any] struct {
	Name         string
	Reader       func(string, *Config) error
	Value        func(Config) string
	Description  string
	Instructions string
}

func (e FromConfigFile) Apply(config *Config) Error {
	data := map[string]string{}
	lErr := yaml.Unmarshal([]byte(e.FileContents), &data)
	if lErr != nil {
		return NewError(
			"Invalid config file (~/"+KLOG_FOLDER+CONFIG_FILE+")",
			"The syntax in the file is not valid YAML.",
			lErr,
		)
	}
	for key, value := range data {
		for _, entry := range CONFIG_FILE_ENTRIES {
			if entry.Name == key {
				rErr := entry.Reader(value, config)
				if rErr != nil {
					return NewError(
						"Invalid config file (~/"+KLOG_FOLDER+CONFIG_FILE+")",
						"The value for `"+key+"` is not valid: "+entry.Instructions,
						rErr,
					)
				}
			}
		}
	}
	return nil
}
