package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Db           Db
		App          App
		JwtSecretKey string
	}

	App struct {
		Host string
		Port int
	}

	// Database
	Db struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
		TimeZone string
	}
)

func GetConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %v", err))
	}

	return Config{
		App: App{
			Host: viper.GetString("app.server.host"),
			Port: viper.GetInt("app.server.port"),
		},
		Db: Db{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.dbname"),
			SSLMode:  viper.GetString("database.sslmode"),
			TimeZone: viper.GetString("database.timezone"),
		},
		JwtSecretKey: viper.GetString("jwt.secret"),
	}
}
