package conf

import "gopkg.in/yaml.v3"

type RoutesConfig struct {
	Routes []struct {
		RoutingKey string `yaml:"routing_key"`
		API        string `yaml:"api"`
	} `yaml:"routes"`
}

func UnmarshalRoutes(data []byte) (RoutesConfig, error) {
	var config RoutesConfig
	err := yaml.Unmarshal(data, &config)
	return config, err
}
