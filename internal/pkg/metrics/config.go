package metrics

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"

	"log"
	"strings"
	"time"
)

var Conf MetricConfig

type MetricConfig struct {
	System    string `env:"METRICS_SYSTEM" envDefault:"nine"`
	SubSystem string `env:"METRICS_SUBSYSTEM" envDefault:"article_dispatcher"`
	HTTP      struct {
		Path         string        `env:"METRICS_HTTP_PATH" envDefault:"metrics"`
		Host         string        `env:"METRICS_HTTP_HOST" envDefault:"7001"`
		ShutdownWait time.Duration `env:"METRICS_HTTP_SERVER_SHUTDOWN_WAIT" envDefault:"5s"`
	}
}

// Register metrics configurations
func (c *MetricConfig) Register() error {
	err := env.Parse(&Conf)
	if err != nil {
		return errors.Wrap(err, "register failed, error loading metrics config")
	}
	return nil
}

// Validate metrics configurations
func (c *MetricConfig) Validate() error {
	if strings.ContainsAny(Conf.System+Conf.SubSystem, "-.") {
		return errors.New("METRICS_SYSTEM and METRICS_SUBSYSTEM variables cannot contain special characters," +
			"_,.")
	}
	if Conf.System == "" || Conf.SubSystem == "" {
		return errors.New("METRICS_SYSTEM or METRICS_SUBSYSTEM variables cannot be empty")
	}
	return nil
}

// Print metrics configurations
func (c *MetricConfig) Print() interface{} {
	defer log.Println("---loading metrics configs---")
	return &Conf
}
