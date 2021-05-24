package clockify

import (
	"github.com/spf13/viper"
)

type Config struct {
	APIKey   string `json:"apiKey"`
	Endpoint string `json:"endpoint"`
}

var cached *Config

// ConfigFromFile reads a configuration file called clockify.yml and returns it as a
// Config instance. If no configuration file is found, nil and no error will be
// returned. The configuration must live in one of the following directories:
//
//	- /etc/timetrace
//	- $HOME/.timetrace
//	- .
//
// In case multiple configuration files are found, the one in the most specific
// or "closest" directory will be preferred.
func ConfigFromFile() (*Config, error) {
	viper.SetConfigName("clockify")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/timetrace/")
	viper.AddConfigPath("$HOME/.timetrace")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	cached = &config

	return cached, nil
}

// GetConfig returns the parsed configuration. The fields of this configuration either
// contain values specified by the user or the zero value of the respective workspaces
// type, e.g. "" for an un-configured string.
//
// Using GetConfig over ConfigFromFile avoids the config file from being parsed each time
// the config is needed.
func GetConfig() *Config {
	if cached != nil {
		return cached
	}

	config, err := ConfigFromFile()
	if err != nil {
		return &Config{}
	}

	return config
}
