package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Port    string `mapstructure:"PORT"`
	DbUrl   string `mapstructure:"DB_URL"`
	ImgPath string `mapstructure:"IMG_PATH"`
	Secret  string `mapstructure:"SECRET"`
}

func LoadConfig() (config Config, err error) {

	viper.SetConfigFile("././.env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&config)
	return
}
