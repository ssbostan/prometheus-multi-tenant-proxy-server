package server

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Global GlobalConfiguration `yaml:"global"`
	Users  []UserConfiguration `yaml:"users"`
}

type GlobalConfiguration struct {
	PrometheusAddress string `yaml:"prometheusAddress"`
}

type UserConfiguration struct {
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Projects []string `yaml:"projects"`
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
