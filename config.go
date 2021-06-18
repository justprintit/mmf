package mmf

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	ClientId     string `yaml:"client_key"`
	ClientSecret string `yaml:"client_secret,omitempty:`
}

func MarshalConfig(c *Config) ([]byte, error) {
	return yaml.Marshal(c)
}

func UnmarshalConfig(b []byte) (*Config, error) {
	var c Config

	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
