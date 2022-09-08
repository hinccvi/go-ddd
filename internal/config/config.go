package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		App  string `mapstructure:"app"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"app"`

	Jwt struct {
		AccessSigningKey  string `mapstructure:"access_signing_key"`
		AccessExpiration  int    `mapstructure:"access_expiration"`
		RefreshSigningKey string `mapstructure:"refresh_signing_key"`
		RefreshExpiration int    `mapstructure:"refresh_expiration"`
	} `mapstructure:"jwt"`

	DB struct {
		Url string `mapstructure:"url"`
	} `mapstructure:"db"`

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

	viper.WatchConfig()

	var errs error
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("-- Config file updated --")
		if err := viper.Unmarshal(conf); err != nil {
			errs = err
		}
	})

	return *conf, errs
}
