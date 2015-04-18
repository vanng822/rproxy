package rproxy

import (
	"os"
	"log"
	"encoding/json"
	"io"
)

type ServerConfig struct {
	ServerName string
	TargetUrl  string
}

type Conf struct {
	Host      string
	Port      int
	ApiEnable bool
	ApiHost   string
	ApiPort   string
	Servers   []*ServerConfig
}

func LoadConfig(filename string) *Conf {
	file, err := os.Open(filename)
	if err != nil {
		log.Panicf("Could not open configuration file, error: %s", err.Error())
	}
	return ParseConfig(file)	
}

func ParseConfig(config io.Reader) *Conf {
	conf := &Conf{}
	decoder := json.NewDecoder(config)
	err := decoder.Decode(conf)
	if err != nil {
		log.Panicf("Could not parse configuration, error: %s", err.Error())
	}
	return conf	
}