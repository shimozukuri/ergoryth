package pgx

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `envconfig:"HOST" default:"localhost"`
	Port     int16         `envconfig:"PORT" default:"5432"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DATABASE" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewConfig() (Config, error) {
	var conf Config

	if err := envconfig.Process("POSTGRES", &conf); err != nil {
		return Config{}, fmt.Errorf(
			"process envconfig: %w",
			err,
		)
	}

	return conf, nil
}

func NewConfigMust() Config {
	conf, err := NewConfig()

	if err != nil {
		err = fmt.Errorf("get Postgres connection pool config: %w", err)
		panic(err)
	}

	return conf
}
