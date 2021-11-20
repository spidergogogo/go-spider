package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type apiKey struct {
	Eth string `yaml:"eth"`
	Bsc string `yaml:"bsc"`
}
type config struct {
	ApiKey apiKey `yaml:"apiKey"`
	Proxy  string `yaml:"proxy"`
}

const configPath = "./cfg.yaml"

var cfg = new(config)

func initCoinConfig() {
	fileData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Error(fmt.Sprintf("读取yaml文件失败:%s", err))
		return
	}
	err = yaml.Unmarshal(fileData, cfg)
	if err != nil {
		log.Error(fmt.Sprintf("解析yaml文件失败:%s", err))
	}
}
