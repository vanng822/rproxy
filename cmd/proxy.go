package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/vanng822/rproxy"
)

func main() {
	var (
		host    string
		port    int
		proxy   string
	)

	flag.StringVar(&host, "h", "127.0.0.1", "Host to listen on")
	flag.IntVar(&port, "p", 80, "Port number to listen on")
	flag.StringVar(&proxy, "r", "", "Proxy host")
	flag.Parse()

	targetUrl, err := url.Parse(proxy)
	
	if err != nil {
		log.Fatal(err)
	}
	
	if targetUrl.Host == "" {
		log.Fatalf("You need to provide a valid proxy url. Provided '%s'", proxy)
	}
	
	rproxy.Server{}
	
	singleHostProxy := httputil.NewSingleHostReverseProxy(targetUrl)
	
	http.Handle("/", singleHostProxy)
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}
