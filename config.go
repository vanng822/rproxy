package rproxy

import (
	"encoding/json"
	"io"
	"log"
	"os"
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
	ApiPort   int
	Servers   []*ServerConfig
}

func DefaultConf() *Conf {
	conf := &Conf{
		Host:      "127.0.0.1",
		Port:      5555,
		ApiEnable: true,
		ApiHost:   "127.0.0.1",
		ApiPort:   5556,
		Servers:   nil,
	}
	return conf
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
