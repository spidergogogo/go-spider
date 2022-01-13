package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type configDaemon struct {
	PidFile string `yaml:"pidFile"`
	LogFile string `yaml:"logFile"`
	WorkDir string `yaml:"workDir"`
}

const daemonConfigPath = "./cfg-daemon.yaml"

func daemonConfig() (*configDaemon, error) {
	cfgDaemon := new(configDaemon)
	fileData, err := ioutil.ReadFile(daemonConfigPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(fileData, cfgDaemon)
	if err != nil {
		return nil, err
	}
	return cfgDaemon, nil
}
