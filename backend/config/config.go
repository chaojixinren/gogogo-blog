package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Port string `mapstructure:"port"`
	} `mapstructure:"app"`
	Database struct {
		DSN          string `mapstructure:"dsn"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
	} `mapstructure:"database"`
	Auth struct {
		JWTSecret     string `mapstructure:"jwt_secret"`
		TokenTTLHours int    `mapstructure:"token_ttl_hours"`
	} `mapstructure:"auth"`
	CORS struct {
		AllowOrigins []string `mapstructure:"allow_origins"`
	} `mapstructure:"cors"`
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	viper.SetDefault("auth.jwt_secret", "change-me")
	viper.SetDefault("auth.token_ttl_hours", 72)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	InitDB()
}
