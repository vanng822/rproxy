package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

	proxyUrl, err := url.Parse(proxy)
	
	if err != nil {
		log.Fatal(err)
	}
	
	if proxyUrl.Host == "" {
		log.Fatalf("You need to provide a valid proxy url. Provided '%s'", proxy)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req = r
			if proxyUrl.Scheme == "" {
				proxyUrl.Scheme = "http"
			}
			req.URL.Scheme = proxyUrl.Scheme
			req.URL.Host = proxyUrl.Host
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}
