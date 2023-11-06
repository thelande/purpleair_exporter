package main

import (
	"errors"
	"os"

	"github.com/go-kit/log/level"
	"gopkg.in/yaml.v3"
)

type PurpleAirConfig struct {
	ApiKey        string   `yaml:"api_key"`
	SensorIndices []int    `yaml:"sensor_indices"`
	Fields        []string `yaml:"fields"`
}

func (c *PurpleAirConfig) ReloadConfig(configFile string) error {
	var config []byte
	var err error

	if configFile != "" {
		config, err = os.ReadFile(configFile)
		if err != nil {
			level.Error(logger).Log("msg", "Error reading config file", "err", err)
			return err
		}
	} else {
		return errors.New("no configuration file specified")
	}

	if err = yaml.Unmarshal(config, c); err != nil {
		return err
	}

	level.Info(logger).Log("msg", "Loaded configuration file", "path", configFile)
	return nil
}
