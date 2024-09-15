package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PostgresURL      string `mapstructure:"POSTGRES_URL"`
	ElasticsearchURL string `mapstructure:"ELASTICSEARCH_URL"`
	ServiceBusConnStr string `mapstructure:"SERVICEBUS_CONNECTION_STRING"`
	ServiceBusQueue   string `mapstructure:"SERVICEBUS_QUEUE_NAME"`
	ServerPort       string `mapstructure:"SERVER_PORT"`
}

func LoadConfig() (config Config, err error) {
	viper.AutomaticEnv()

	viper.SetDefault("POSTGRES_URL", "postgres://user:password@localhost:5432/dbname")
	viper.SetDefault("ELASTICSEARCH_URL", "http://localhost:9200")
	viper.SetDefault("SERVICEBUS_CONNECTION_STRING", "")
	viper.SetDefault("SERVICEBUS_QUEUE_NAME", "events")
	viper.SetDefault("SERVER_PORT", "8080")

	err = viper.Unmarshal(&config)
	return
}