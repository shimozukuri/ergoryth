package logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env string `envconfig:"ENV" required:"true"`
}

func NewConfig() (Config, error) {
	var conf Config

	if err := envconfig.Process("LOGGER", &conf); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return conf, nil
}

func NewConfigMust() Config {
	conf, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Logger config: %w", err)
		panic(err)
	}

	return conf
}
