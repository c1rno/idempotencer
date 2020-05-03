package config

import (
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/logging"
	"github.com/c1rno/idempotencer/pkg/metrics"
	"github.com/c1rno/idempotencer/pkg/persistence"
	"github.com/c1rno/idempotencer/pkg/queue"
	"github.com/c1rno/idempotencer/pkg/upstream"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel      int
	Metrics       metrics.Config
	Upstream      upstream.Config
	QueueConsumer queue.ConsumerConfig
	QueueProducer queue.ProducerConfig
	QueueBroker   queue.BrokerConfig
	Persistence   persistence.Config
}

func NewConfig(l logging.Logger) (c Config, err error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("default-config")
	viper.AddConfigPath("/etc/idempotencer")
	if err = viper.ReadInConfig(); err != nil {
		l.Error(errors.NewError(errors.ConfigReadingError, err).String(), nil)
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&c)
	if err != nil {
		err = errors.NewError(errors.ConfigUnmarshallError, err)
	}
	return c, err
}
