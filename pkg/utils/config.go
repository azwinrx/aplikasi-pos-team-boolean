package utils

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Config Configuration

type Configuration struct {
	AppName     string
	Port        string
	Env         string
	Debug       bool
	Limit       int
	PathLogging string
	JWTSecret   string
	DB          DatabaseCofig
	SMTP        SMTPConfig
}

type DatabaseCofig struct {
	Name     string
	Username string
	Password string
	Host     string
	Port     string
	MaxConn  int32
}

type SMTPConfig struct {
	Host     string
	Port     string
	Email    string
	Password string
	FromName string
}

func ReadConfiguration() (Configuration, error) {
	// get config from env file
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return Configuration{}, err
	}

	// get config from os variable
	viper.AutomaticEnv()

	// get config from flag
	pflag.Int("port-app", 0, "port for app golang")
	viper.BindPFlags(pflag.CommandLine)

	Config = Configuration{
		AppName:     viper.GetString("APP_NAME"),
		Port:        viper.GetString("PORT"),
		Env:         viper.GetString("ENV"),
		Debug:       viper.GetBool("DEBUG"),
		Limit:       viper.GetInt("LIMIT"),
		PathLogging: viper.GetString("PATH_LOGGING"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		DB: DatabaseCofig{
			Name:     viper.GetString("DATABASE_NAME"),
			Username: viper.GetString("DATABASE_USERNAME"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			Host:     viper.GetString("DATABASE_HOST"),
			Port:     viper.GetString("DATABASE_PORT"),
			MaxConn:  viper.GetInt32("DATABASE_MAX_CONN"),
		},
		SMTP: SMTPConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetString("SMTP_PORT"),
			Email:    viper.GetString("SMTP_EMAIL"),
			Password: viper.GetString("SMTP_PASSWORD"),
			FromName: viper.GetString("SMTP_FROM_NAME"),
		},
	}
	return Config, nil

}
