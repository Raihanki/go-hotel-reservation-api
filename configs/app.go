package configs

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Config struct {
	APP_NAME         string
	APP_PORT         string
	APP_ENV          string
	APP_URL          string
	APP_FRONTEND_URL string

	DB_CONNECTION string
	DB_HOST       string
	DB_PORT       string
	DB_DATABASE   string
	DB_USERNAME   string
	DB_PASSWORD   string

	JWT_SECRET string

	STRIPE_KEY string
}

var ENV *Config

func LoadConfigApp() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")

	errReadConfig := viper.ReadInConfig()
	if errReadConfig != nil {
		log.Fatal("Error reading config file")
	}

	errUnmarshal := viper.Unmarshal(&ENV)
	if errUnmarshal != nil {
		log.Fatal("Error unmarshal config file")
	}

	log.Info("Config successfully loaded")
}
