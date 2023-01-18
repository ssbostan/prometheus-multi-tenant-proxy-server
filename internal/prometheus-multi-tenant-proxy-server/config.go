package server

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Global GlobalConfig `yaml:"global"`
	Users  []UserConfig `yaml:"users"`
}

type GlobalConfig struct {
	ListenAddress       string `yaml:"listenAddress"`
	PrometheusAddress   string `yaml:"prometheusAddress"`
	AccessRequestHeader string `yaml:"accessRequestHeader"`
	AccessTargetLabel   string `yaml:"accessTargetLabel"`
}

type UserConfig struct {
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Accesses []string `yaml:"accesses"`
}

func parseConfigFile(configFile string) Config {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
