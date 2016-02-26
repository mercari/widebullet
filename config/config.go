package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

const (
	DefaultPort                = "29300"
	DefaultLogLevel            = "error"
	DefaultTimeout             = 5
	DefaultMaxIdleConnsPerHost = 100
)

type Config struct {
	Port                string
	LogLevel            string
	Timeout             int
	MaxIdleConnsPerHost int
	DisableCompression  bool
	Endpoints           []EndPoint
}

type EndPoint struct {
	Name            string
	Ep              string
	ProxySetHeaders [][]string
}

func Load(confPath string) (Config, error) {
	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return Config{}, err
	}

	var config Config

	if err := toml.Unmarshal(bytes, &config); err != nil {
		return config, err
	}

	if config.Port == "" {
		config.Port = DefaultPort
	}

	if config.LogLevel == "" {
		config.LogLevel = DefaultLogLevel
	}

	if config.Timeout <= 0 {
		config.Timeout = DefaultTimeout
	}

	if config.MaxIdleConnsPerHost <= 0 {
		config.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	}

	if len(config.Endpoints) == 0 {
		return config, errors.New("empty Endpoints")
	}

	for _, ep := range config.Endpoints {
		if ep.Name == "" {
			return config, errors.New("empty Endpoint name")
		}
		if ep.Ep == "" {
			return config, errors.New("empty Endpoint URL")
		}
	}

	return config, nil
}

func FindEp(conf Config, name string) (EndPoint, error) {
	for _, ep := range conf.Endpoints {
		if ep.Name == name {
			return ep, nil
		}
	}

	return EndPoint{}, fmt.Errorf("ep:%s is not found", name)
}
