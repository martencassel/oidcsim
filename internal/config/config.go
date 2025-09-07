package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	OIDC struct {
		Issuer  string `yaml:"issuer"`
		Signing struct {
			PrivateKeyFile string `yaml:"privateKeyFile"`
		} `yaml:"signing"`
	} `yaml:"oidc"`

	Routes struct {
		Token     string `yaml:"token"`
		Authorize string `yaml:"authorize"`
	} `yaml:"routes"`
}

func Load(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
