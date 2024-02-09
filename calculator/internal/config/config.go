package config

import (
	"github.com/spf13/viper"
)

// Config структура для работы с конфигурацией.
type Config struct {
	ConnString string `mapstructure:"CONN_STRING"`
	DriverName string `mapstructure:"DRIVER_NAME"`
}

// MustConfig парсит конфиг из .env файла, в случае ошибки паникует.
func MustConfig() *Config {
	var config Config
	viper.AddConfigPath("./internal/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic("failed to read config: " + err.Error())
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err.Error())
	}
	return &config
}
