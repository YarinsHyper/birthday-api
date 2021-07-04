package util

import (
	"github.com/spf13/viper"
)

// Config is used to declare every variable in
// app.env file in order to use it in go files
type Config struct {
	Port       string `mapstructure:"PORT"`
	RouterPort string `mapstructure:"ROUTER_PORT"`
	Address    string `mapstructure:"ADDRESS"`
}

// LoadConfig loads all variables from app.env file
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
