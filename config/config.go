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

	return viper.ReadInConfig()
}

func setDefaults() {
	viper.SetDefault("port", 8080)
	viper.SetDefault("readTimeout", 30*time.Millisecond)
	viper.SetDefault("writeTimeout", 30*time.Millisecond)
	viper.SetDefault("idleTimeout", 1*time.Second)
}

func ListenAddress() string {
	return fmt.Sprintf(":%d", viper.GetInt("port"))
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
