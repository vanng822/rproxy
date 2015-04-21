package main

import (
	"flag"
	"github.com/vanng822/rproxy"
)

func main() {
	var (
		host    string
		port    int
		apiHost string
		apiPort int
		config  string
		conf    *rproxy.Conf
	)

	flag.StringVar(&host, "h", "", "Host to listen on")
	flag.IntVar(&port, "p", -1, "Port number to listen on")
	flag.StringVar(&apiHost, "ah", "", "Host for server admin api")
	flag.IntVar(&apiPort, "ap", -1, "Port for server admin api")
	flag.StringVar(&config, "c", "", "Configuration file")
	flag.Parse()
	
	if config != "" {
		conf = rproxy.LoadConfig(config)
	} else {
		conf = rproxy.DefaultConf()
	}
	
	if host != "" {
		conf.Host = host
	}
	
	if port != -1 {
		conf.Port = port
	}
	
	if apiHost != "" {
		conf.ApiHost = apiHost
	}
	
	if apiPort != -1 {
		conf.ApiPort = apiPort
	}
	
	proxyServer := rproxy.NewProxy(conf)
	proxyServer.Start()
}
