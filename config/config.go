package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

var config *Config

//Config all configurable options for the parser and routingplaner
type Config struct {
	OsmLocation string `yaml:"osmLocation"`
	OsmFilename string `yaml:"osmFilename"`
	OsmParse    int    `yaml:"osmParse"`

	OutputType     string `yaml:"outputType"`
	OutputFilename string `yaml:"outputFilename"`

	LogLevel int `yaml:"logLevel"`

	GridXSize int `yaml:"gridXSize"`
	GridYSize int `yaml:"gridYSize"`

	InfoFilename string `yaml:"infoFilename"`
}

//LoadConfig .
func LoadConfig(filename string) *Config {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {

		log.Fatal(err)
	}

	config = &Config{}

	err2 := yaml.Unmarshal(fileData, config)
	if err != nil {

		log.Fatal(err2)

	}

	return config
}

//GetConfig .
func GetConfig() *Config {
	if config == nil {
		log.Fatal("Load config first.")
	}
	return config
}
