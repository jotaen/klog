package app

import "github.com/pelletier/go-toml"

type Config struct {
	DateWithDashes  bool
	Time24HourClock bool
	TrimTime        bool
}

func NewDefaultConfig() Config {
	return Config{
		DateWithDashes:  true,
		Time24HourClock: true,
		TrimTime:        false,
	}
}

func NewConfigFromToml(tomlText string) (Config, error) {
	userConfig, err := toml.Load(tomlText)
	defaultConfig := NewDefaultConfig()
	if err != nil {
		return defaultConfig, err
	}

	if userConfig.Get("format.date").(string) == "YYYY/MM/DD" {
		defaultConfig.DateWithDashes = true
	}

	defaultConfig.Time24HourClock = userConfig.GetDefault("format.clock_24_hours", defaultConfig.Time24HourClock).(bool)
	defaultConfig.TrimTime = userConfig.GetDefault("format.trim_time", defaultConfig.TrimTime).(bool)

	return defaultConfig, nil
}

func DefaultConfigAsToml() string {
	return `########## klog Configuration ##########
# klog time tracking
# https://www.github.com/jotaen/klog
########################################

[format]
##### How dates are printed out.
##### "YYYY-MM-DD"     E.g. 2020-04-26   [default]
##### "YYYY/MM/DD"     E.g. 2020/04/26
# date = "YYYY-MM-DD"

##### Whether to use the 24-hour or the 12-hour clock.
##### true             E.g. 19:30        [default]
##### false            E.g. 7:30pm
# clock_24_hours = true

##### ===== Trim time =====
##### Whether to trim leading zeros in the hour part of times.
##### true             E.g. 6:05         [default]
##### false            E.g. 06:05
# trim_time = true
`
}
