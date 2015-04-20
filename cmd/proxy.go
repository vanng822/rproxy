package main

import (
	"flag"
	"fmt"
	"github.com/vanng822/accesslog"
	"github.com/vanng822/r2router"
	"github.com/vanng822/recovery"
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

	logger := accesslog.New()
	rec := recovery.NewRecovery()

	seefor := r2router.NewSeeforRouter()
	seefor.Before(rec.Handler)
	seefor.Before(logger.Handler)

	seefor.Group("/_server", func(r *r2router.GroupRouter) {
		r.Post("/backend", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
			err, severConfig := proxyServer.ParseServerConfig(req)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid server config, error: %s", err.Error()), http.StatusBadRequest)
				return
			}
			err = proxyServer.Register(severConfig.ServerName, severConfig.TargetUrl)
			if err != nil {
				http.Error(w,
					fmt.Sprintf("It was problem when adding new server, serverName: '%s', targetUrl: '%s', error: '%s'",
						severConfig.ServerName, severConfig.TargetUrl, err.Error()),
					http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))

		})
		// delete backend node
		r.Delete("/backend", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
			err, severConfig := proxyServer.ParseServerConfig(req)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid server config, error: %s", err.Error()), http.StatusBadRequest)
				return
			}
			err = proxyServer.Unregister(severConfig.ServerName, severConfig.TargetUrl)
			if err != nil {
				http.Error(w,
					fmt.Sprintf("It was problem when removing server, serverName: '%s', targetUrl: '%s', error: '%s'",
						severConfig.ServerName, severConfig.TargetUrl, err.Error()),
					http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))

		})
		// delete server
		r.Delete("/", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
			req.ParseForm()
			serverName := req.Form.Get("serverName")
			if serverName == "" {
				http.Error(w, fmt.Sprintf("serverName is required"), http.StatusBadRequest)
				return
			}
			err := proxyServer.RemoveServer(serverName)
			if err != nil {
				http.Error(w,
					fmt.Sprintf("It was problem when removing server, serverName: '%s', error: '%s'",
						severConfig.ServerName, err.Error()),
					http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	})

	http.Handle("/", rec.Handler(logger.Handler(proxyServer)))
	go http.ListenAndServe(fmt.Sprintf("%s:%d", apiHost, apiPort), seefor)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}
