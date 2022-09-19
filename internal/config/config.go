package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"app"`

	Jwt struct {
		AccessSigningKey  string `mapstructure:"access_signing_key"`
		AccessExpiration  int    `mapstructure:"access_expiration"`
		RefreshSigningKey string `mapstructure:"refresh_signing_key"`
		RefreshExpiration int    `mapstructure:"refresh_expiration"`
	} `mapstructure:"jwt"`

	Dsn string `mapstructure:"dsn"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
		PoolSize int    `mapstructure:"pool_size"`
	} `mapstructure:"redis"`
}

func Load(mode string) (Config, error) {
	file := ""

	switch mode {
	case "local":
		file = "local"
	case "dev":
		file = "dev"
	case "qa":
		file = "qa"
	case "prod":
		file = "prod"
	}

	viper.SetConfigName(file)
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../../config/")

	conf := new(Config)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return *conf, err
		} else {
			return *conf, err
		}
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return *conf, err
	}

	return *conf, nil
}
