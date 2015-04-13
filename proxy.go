package rproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	// Matching against req.Host
	name    string
	backend *Backend
}

type Proxy struct {
	servers map[string]*Server
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for name, server := range p.servers {
		if name == req.Host {
			server.backend.ServeHTTP(rw, req)
			return
		}
	}
}

func (p *Proxy) Register(serverName, targetUrl string) error {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return err
	}
	b := &Backend{}
	b.nodes[targetUrl] = httputil.NewSingleHostReverseProxy(target)
	p.servers[serverName] = &Server{
		name:    serverName,
		backend: b,
	}
	return nil
}

func (p *Proxy) Unregister(serverName, targetUrl string) error {
	s, ok := p.servers[serverName]

	if !ok {
		return fmt.Errorf("No server by name %s", serverName)
	}

	if _, ok := s.backend.nodes[targetUrl]; !ok {
		return fmt.Errorf("No target %s exists for servername %s", targetUrl, serverName)
	}

	delete(s.backend.nodes, targetUrl)

	return nil
}
