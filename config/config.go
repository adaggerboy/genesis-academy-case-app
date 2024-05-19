package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	model "github.com/adaggerboy/genesis-academy-case-app/models/config"
)

var GlobalConfig model.Config

func Load(filename string) (config model.Config, err error) {
	err = nil
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("config loading: %s", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		err = fmt.Errorf("config parsing: %s", err)
	}
	return
}
