package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

var config *Config

type Config struct {
	OsmLocation string `yaml:"osmLocation"`
	OsmFilename string `yaml:"osmFilename"`
	OsmParse    int    `yaml:"osmParse"`

	OutputType     string `yaml:"outputType"`
	OutputFilename string `yaml:"outputFilename"`
}

func LoadConfig(filename string) *Config {
	file_data, err := ioutil.ReadFile(filename)
	if err != nil {

		log.Fatal(err)
	}

	config = &Config{}

	err2 := yaml.Unmarshal(file_data, config)
	if err != nil {

		log.Fatal(err2)

	}

	return config
}

func GetConfig() *Config {
	if config == nil {
		log.Fatal("Load config first.")
	}
	return config
}
