package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func Init() error {
	viper.SetEnvPrefix("chinnaswamy")
	viper.AutomaticEnv()
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("redirectServerPort", 8080)
	viper.SetDefault("shortenServerPort", 8090)
	viper.SetDefault("readTimeout", 30*time.Millisecond)
	viper.SetDefault("writeTimeout", 30*time.Millisecond)
	viper.SetDefault("idleTimeout", 1*time.Second)
}

func RedirectListenAddress() string {
	return fmt.Sprintf(":%d", viper.GetInt("redirectServerPort"))
}

func ShortenListenAddress() string {
	return fmt.Sprintf(":%d", viper.GetInt("shortenServerPort"))
}

func ReadTimeout() time.Duration {
	return viper.GetDuration("readTimeout")
}

func WriteTimeout() time.Duration {
	return viper.GetDuration("writeTimeout")
}

func IdleTimeout() time.Duration {
	return viper.GetDuration("idleTimeout")
}

func DatabaseAddresses() []string {
	return viper.GetStringSlice("databaseAddresses")
}
