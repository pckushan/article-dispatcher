package log

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"log"
)

var Config LoggerConfig

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"TRACE"`
}

// Register log configurations
func (l *LoggerConfig) Register() error {
	err := env.Parse(&Config)
	if err != nil {
		return errors.Wrap(err, "register failed, error parsing logger config")
	}
	return nil
}

// Validate log configurations
func (l *LoggerConfig) Validate() error {
	return nil
}

// Print log configurations
func (l *LoggerConfig) Print() interface{} {
	defer log.Println("---loading logger configs---")
	return Config
}
