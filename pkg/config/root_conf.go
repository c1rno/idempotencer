package config

import (
	"github.com/c1rno/idempotencer/pkg/consumer"
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/logging"
	"github.com/c1rno/idempotencer/pkg/persistence"
	"github.com/c1rno/idempotencer/pkg/producer"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel    int
	Consumer    consumer.Config
	Producer    producer.Config
	Persistence persistence.Config
}

func NewConfig(l logging.Logger) (c Config, err error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("default-config")
	viper.AddConfigPath("/etc/idempotencer")
	if err = viper.ReadInConfig(); err != nil {
		l.Error(errors.NewError(errors.ConfigReadingError, err).String(), nil)
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		err = errors.NewError(errors.ConfigUnmarshallError, err)
	}
	return c, err
}
