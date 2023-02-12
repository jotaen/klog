package app

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
}

type ConfigDeterminer interface {
	Apply(*Config)
}

func NewConfig(c1 ConfigFromStaticValues, c2 ConfigFromEnvVars) Config {
	config := NewDefaultConfig()
	c1.Apply(&config)
	c2.Apply(&config)
	return config
}

func NewDefaultConfig() Config {
	return Config{
		IsDebug:    newDefaultParam(false),
		Editor:     newDefaultParam(""),
		NoColour:   newDefaultParam(false),
		CpuKernels: newDefaultParam(1),
	}
}

type Param[T any] struct {
	actualValue  T
	defaultValue T
	isExplicit   bool
}

func newDefaultParam[T any](value T) Param[T] {
	return Param[T]{
		actualValue:  value,
		defaultValue: value,
		isExplicit:   false,
	}
}

func (p Param[T]) Value() T {
	return p.actualValue
}

func (p Param[T]) Default() T {
	return p.defaultValue
}

func (p Param[T]) WasExplicitlySet() bool {
	return p.isExplicit
}

func (p *Param[T]) set(value T) {
	p.actualValue = value
	p.isExplicit = true
}

type ConfigFromStaticValues struct {
	NumCpus int
}

func (e ConfigFromStaticValues) Apply(config *Config) {
	config.CpuKernels.set(e.NumCpus)
}

type ConfigFromEnvVars struct {
	GetVar func(string) string
}

func (e ConfigFromEnvVars) Apply(config *Config) {
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
}
