package app

import "github.com/pelletier/go-toml"

type Config struct {
	DateWithDashes  bool
	Time24HourClock bool
}

func NewDefaultConfig() Config {
	return Config{
		DateWithDashes:  true,
		Time24HourClock: true,
	}
}

func NewConfigFromToml(tomlText string) (Config, error) {
	userConfig, err := toml.Load(tomlText)
	defaultConfig := NewDefaultConfig()
	if err != nil {
		return defaultConfig, err
	}

	if userConfig.Get("date_format").(string) == "YYYY/MM/DD" {
		defaultConfig.DateWithDashes = false
	}

	if userConfig.Get("clock_convention").(string) == "12_hours" {
		defaultConfig.Time24HourClock = false
	}

	return defaultConfig, nil
}

func DefaultConfigAsToml() string {
	return `########## klog Configuration ##########
# klog time tracking
# https://www.github.com/jotaen/klog
########################################

##### How dates are printed out:
##### "YYYY-MM-DD"     E.g. 2020-04-26   [default]
##### "YYYY/MM/DD"     E.g. 2020/04/26
# date_format = "YYYY-MM-DD"

##### Which clock convention to use:
##### "24_hours"       E.g. 19:30        [default]
##### "12_hours"       E.g. 7:30pm
# clock_convention = "24_hours"
`
}
