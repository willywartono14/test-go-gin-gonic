package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	config *Config
)

func Init(configFile string) error {

	out, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(out, &config)
	if err != nil {
		return err
	}

	return nil
}

func Get() *Config {
	return config
}
