package main

import (
	"flag"
	"fmt"
	"github.com/vanng822/rproxy"
	"github.com/vanng822/accesslog"
	"github.com/vanng822/recovery"
	"log"
	"net/http"
)

func main() {
	var (
		host    string
		port    int
		apiHost string
		apiPort int
	)

	flag.StringVar(&host, "h", "127.0.0.1", "Host to listen on")
	flag.IntVar(&port, "p", 80, "Port number to listen on")
	flag.StringVar(&apiHost, "ah", "127.0.0.1", "Host for server admin api")
	flag.IntVar(&apiPort, "ap", 8080, "Port for server admin api")
	flag.Parse()

	logger := accesslog.New()
	rec := recovery.NewRecovery()
	
	proxyServer := rproxy.NewProxy()
	
	api := proxyServer.AdminAPI()
	
	http.Handle("/", rec.Handler(logger.Handler(proxyServer)))
	go http.ListenAndServe(fmt.Sprintf("%s:%d", apiHost, apiPort), api)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}
