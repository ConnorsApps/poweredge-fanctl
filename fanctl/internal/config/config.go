package config

import (
	"os"

	"github.com/ConnorsApps/poweredge-fanctl/pkg/ipmitool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func SetLogLevel(level string) {
	if l, err := zerolog.ParseLevel(level); err != nil {
		log.Error().Err(err).Str("level", level).Msg("Error parsing log level")
	} else {
		zerolog.SetGlobalLevel(l)
	}
}

type Config struct {
	IDRAC    *ipmitool.Config `yaml:"idrac"`
	LogLevel string           `yaml:"logLevel"`
}

func MustRead(path string) *Config {
	var c Config
	if path == "" {
		path = "config.yaml"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panic().Err(err).Str("path", path).Msg("Error reading config file")
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		log.Panic().Err(err).Msg("Error unmarshalling config file")
	}
	if c.LogLevel != "" {
		SetLogLevel(c.LogLevel)
	}
	return &c
}
