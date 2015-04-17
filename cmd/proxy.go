package main

import (
	"flag"
	"fmt"
	"github.com/vanng822/r2router"
	"github.com/vanng822/rproxy"
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
	
	proxyServer := rproxy.NewProxy()

	seefor := r2router.NewSeeforRouter()

	seefor.Group("/_server", func(r *r2router.GroupRouter) {
		r.Post("/add", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
			log.Println("_server/add")
			err, serverName, targetUrl := proxyServer.ParseServerConfig(req)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid server config, error: %s", err.Error()), http.StatusBadRequest)
				return
			}
			err = proxyServer.Register(serverName, targetUrl)
			if err != nil {
				http.Error(w,
					fmt.Sprintf("It was problem when adding new server, serverName: '%s', targetUrl: '%s', error: '%s'",
						serverName, targetUrl, err.Error()),
					http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))

		})

		r.Delete("/remove", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
			log.Println("_server/remove")
			err, serverName, targetUrl := proxyServer.ParseServerConfig(req)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid server config, error: %s", err.Error()), http.StatusBadRequest)
				return
			}
			err = proxyServer.Unregister(serverName, targetUrl)
			if err != nil {
				http.Error(w,
					fmt.Sprintf("It was problem when removing server, serverName: '%s', targetUrl: '%s', error: '%s'",
						serverName, targetUrl, err.Error()),
					http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))

		})
	})

	http.Handle("/", proxyServer)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", apiHost, apiPort), seefor)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}
