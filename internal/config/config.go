package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	AppConfig struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"app"`

	JwtConfig struct {
		JWTSigningKey string `mapstructure:"jwt_signing_key"`
		JWTExpiration string `mapstructure:"jwt_expiration"`
	} `mapstructure:"jwt"`

	DBConfig struct {
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		Port     int    `mapstructure:"port"`
	} `mapstructure:"db"`

	RedisConfig struct {
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
