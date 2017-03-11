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
	DefaultIdleConnTimeout     = 30
	DefaultProxyReadTimeout    = 60
	DefaultShutdownTimeout     = 10
)

type Config struct {
	Port                string
	LogLevel            string
	Timeout             int
	MaxIdleConnsPerHost int
	DisableCompression  bool
	IdleConnTimeout     int
	ProxyReadTimeout    int
	ShutdownTimeout     int
	Endpoints           []EndPoint
}

type EndPoint struct {
	Name             string
	Ep               string
	ProxySetHeaders  [][]string
	ProxyPassHeaders [][]string
}

func LoadBytes(bytes []byte) (Config, error) {
	var config Config
	if err := toml.Unmarshal(bytes, &config); err != nil {
		return config, err
	}
	return config, nil
}

func Load(confPath string) (Config, error) {
	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return Config{}, err
	}

	config, err := LoadBytes(bytes)
	if err != nil {
		return Config{}, err
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

	if config.IdleConnTimeout <= 0 {
		config.IdleConnTimeout = DefaultIdleConnTimeout
	}

	if config.ProxyReadTimeout <= 0 {
		config.ProxyReadTimeout = DefaultProxyReadTimeout
	}

	if config.ShutdownTimeout <= 0 {
		config.ShutdownTimeout = DefaultShutdownTimeout
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
