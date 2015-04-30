package rproxy

import (
	"fmt"
	"github.com/unrolled/render"
	"github.com/vanng822/accesslog"
	"github.com/vanng822/r2router"
	"github.com/vanng822/recovery"
	"net/http"
)

func (p *Proxy) AdminAPI() *r2router.Seefor {

	logger := accesslog.New()
	rec := recovery.NewRecovery()

	renderer := render.New()

	seefor := r2router.NewSeeforRouter()
	seefor.Before(rec.Handler)
	seefor.Before(logger.Handler)

	seefor.Group("/_server", func(r *r2router.GroupRouter) {
		r.Post("/backend", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
			err, severConfig := p.ParseServerConfig(req)
			if err != nil {
				renderer.JSON(w,
					http.StatusBadRequest,
					r2router.M{
						"status": "ERROR",
						"msg":    fmt.Sprintf("Invalid server config, error: %s", err.Error()),
					})
				return
			}
			err = p.Register(severConfig.ServerName, severConfig.TargetUrl)
			if err != nil {
				renderer.JSON(w,
					http.StatusInternalServerError,
					r2router.M{
						"status": "ERROR",
						"msg": fmt.Sprintf("It was problem when adding new server, serverName: '%s', targetUrl: '%s', error: '%s'",
							severConfig.ServerName, severConfig.TargetUrl, err.Error()),
					})
				return
			}
			renderer.JSON(w, http.StatusOK, r2router.M{"status": "OK"})
		})
		// delete backend node
		r.Delete("/backend", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
			err, severConfig := p.ParseServerConfig(req)
			if err != nil {
				renderer.JSON(w,
					http.StatusBadRequest,
					r2router.M{
						"status": "ERROR",
						"msg":    fmt.Sprintf("Invalid server config, error: %s", err.Error()),
					})
				return
			}
			err = p.Unregister(severConfig.ServerName, severConfig.TargetUrl)
			if err != nil {
				renderer.JSON(w,
					http.StatusInternalServerError,
					r2router.M{
						"status": "ERROR",
						"msg": fmt.Sprintf("It was problem when removing server, serverName: '%s', targetUrl: '%s', error: '%s'",
							severConfig.ServerName, severConfig.TargetUrl, err.Error()),
					})
				return
			}
			renderer.JSON(w, http.StatusOK, r2router.M{"status": "OK"})

		})
		// delete server
		r.Delete("/", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
			req.ParseForm()
			serverName := req.Form.Get("serverName")
			if serverName == "" {
				renderer.JSON(w,
					http.StatusBadRequest,
					r2router.M{
						"status": "ERROR",
						"msg":    fmt.Sprintf("serverName is required"),
					})
				return
			}
			err := p.RemoveServer(serverName)
			if err != nil {
				renderer.JSON(w,
					http.StatusInternalServerError,
					r2router.M{
						"status": "ERROR",
						"mgs": fmt.Sprintf("It was problem when removing server, serverName: '%s', error: '%s'",
							serverName, err.Error()),
					})
				return
			}
			renderer.JSON(w, http.StatusOK, r2router.M{"status": "OK"})
		})
		// list all node for a servername
		r.Get("/backend", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
			req.ParseForm()
			serverName := req.Form.Get("serverName")
			if serverName == "" {
				renderer.JSON(w,
					http.StatusBadRequest,
					r2router.M{
						"status": "ERROR",
						"msg":    fmt.Sprintf("serverName is required"),
					})
				return
			}
			server, ok := p.servers[serverName]
			if !ok {
				renderer.JSON(w,
					http.StatusBadRequest,
					r2router.M{
						"status": "ERROR",
						"msg":    fmt.Sprintf("Could not find server by name: %s", serverName),
					})
				return
			}
			data := r2router.M{}
			nodes := make([]string, len(server.backend.nodes))
			for i, node := range server.backend.nodes {
				nodes[i] = node.targetUrl
			}
			data[serverName] = nodes
			renderer.JSON(w, http.StatusOK, r2router.M{"status": "OK", "result": data})
		})
		// list all servernames
		r.Get("/", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
			data := r2router.M{}
			for serverName, server := range p.servers {
				nodes := make([]string, len(server.backend.nodes))
				for i, node := range server.backend.nodes {
					nodes[i] = node.targetUrl
				}
				data[serverName] = nodes
			}
			renderer.JSON(w, http.StatusOK, r2router.M{"status": "OK", "result": data})
		})
	})

	return seefor
}
