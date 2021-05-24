package clockify

import "github.com/dominikbraun/timetrace/config"

type Config struct {
	config.Config

	Clockify struct {
		APIKey   string `json:"apiKey"`
		Endpoint string `json:"endpoint"`
	}
}
