package configs

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/configs"
)

type Config struct {
	configs.BaseConfig `yaml:",inline"`
	MQTT               *MQTT `yaml:"mqtt"`
}

type MQTT struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	rawContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(rawContent, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
