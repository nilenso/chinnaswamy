package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"nilenso.com/chinnaswamy/log"
	"time"
)

type config struct {
	Port          uint16        `yaml:"port"`
	WriteTimeout_ time.Duration `yaml:"writeTimeout"`
	ReadTimeout_  time.Duration `yaml:"readTimeout"`
	IdleTimeout_  time.Duration `yaml:"idleTimeout"`
}

type Config interface {
	ListenAddress() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	IdleTimeout() time.Duration
}

func Load(path string) (Config, error) {
	var cfg config
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Cannot open config file: %s", err.Error(),
			"configPath", path,
		)
		return nil, err
	}
	decodeErr := yaml.Unmarshal(contents, &cfg)
	if decodeErr != nil {
		log.Errorf("Cannot decode config file: %s", err.Error(),
			"configPath", path,
		)
	}
	return &cfg, nil
}

func (cfg *config) ListenAddress() string {
	return fmt.Sprintf(":%d", cfg.Port)
}

func (cfg *config) ReadTimeout() time.Duration {
	return cfg.ReadTimeout_
}

func (cfg *config) WriteTimeout() time.Duration {
	return cfg.WriteTimeout_
}

func (cfg *config) IdleTimeout() time.Duration {
	return cfg.IdleTimeout_
}
