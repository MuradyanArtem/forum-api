package apiserver

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config - server configuration datastructure
type Config struct {
	BindAddr    string `yaml:"bind_addr"`
	LogLevel    string `yaml:"log_level"`
	DatabaseURL string `yaml:"database_url"`
	SessionKey  string `yaml:"session_key"`
}

// ParseConfig read configuration from file to struct
func ParseConfig(configPath string, serverConfig *Config) error {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &serverConfig)
	if err != nil {
		return err
	}
	return nil
}
