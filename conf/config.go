package conf

import (
	"fmt"
	"github.com/caarlos0/env"
	"io/ioutil"
)

type Config struct {
	ConfigRoutingKey string `env:"CONFIG_ROUTING_KEY" envDefault:"conf/config.yml"`
	Port             string `env:"PORT" envDefault:"8081"`
	ERPDomain        string `env:"ERP_DOMAIN" envDefault:"http://localhost:8000"`
	ContentType      string `env:"CONTENT_TYPE" envDefault:"application/json"`
}

var cfg Config

func SetEnv() {
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Failed to read environment variables: %v", err)
		return
	}
}

func GetEnv() Config {
	return cfg
}

func GetConfigRoutingKey() (RoutesConfig, error) {
	routesConfig := RoutesConfig{}
	data, err := ioutil.ReadFile(cfg.ConfigRoutingKey)
	if err != nil {
		return routesConfig, err
	}

	routesConfig, err = UnmarshalRoutes(data)
	if err != nil {
		return routesConfig, err
	}

	return routesConfig, nil
}
