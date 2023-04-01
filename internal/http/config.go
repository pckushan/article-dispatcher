package http

import (
	"github.com/caarlos0/env/v6"

	"log"
	"time"
)

var Config RouterConfig

type RouterConfig struct {
	Host     string `env:"HTTP_SERVER_HOST" envDefault:"8888"`
	Timeouts struct {
		Read          time.Duration `env:"HTTP_SERVER_READ_TIMEOUT" envDefault:"10s"`
		Write         time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT" envDefault:"10s"`
		Idle          time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" envDefault:"10s"`
		ShoutDownWait time.Duration `env:"HTTP_SERVER_SHOUT_DOWN_WAIT" envDefault:"5s"`
	}
}

// Register router configurations
func (r *RouterConfig) Register() error {
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal("register failed , error parsing http router config")
	}

	return nil
}

// Validate router configurations
func (r *RouterConfig) Validate() error {
	if Config.Host == "" {
		log.Fatal("application http port cannot be empty")
	}
	return nil
}

// Print router configurations
func (r *RouterConfig) Print() interface{} {
	defer log.Println("---loading router configs---")
	return &Config
}
