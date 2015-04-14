package rproxy

import (
	"fmt"
	"net/http"
	//"net/http/httputil"
	//"net/url"
)

type Server struct {
	// Matching against req.Host
	name    string
	backend *Backend
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	s.backend.ServeHTTP(rw, req)
}

type Proxy struct {
	servers map[string]*Server
}

func NewProxy() *Proxy {
	p := &Proxy{}
	p.servers = make(map[string]*Server)
	return p
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if server, found := p.servers[req.Host]; found {
		server.ServeHTTP(rw, req)
		return
	}
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (p *Proxy) Register(serverName, targetUrl string) error {
	var s *Server
	if _, ok := p.servers[serverName]; ok {
		s = p.servers[serverName]
	} else {
		s = &Server{
			name: serverName,
			backend: NewBackend(),
		}
	}
	
	err := s.backend.addNode(targetUrl)
	
	if err != nil {
		return err
	}
	
	p.servers[serverName] = s
	
	return nil
}

func (p *Proxy) Unregister(serverName, targetUrl string) error {
	s, ok := p.servers[serverName]

	if !ok {
		return fmt.Errorf("No server by name %s", serverName)
	}

	err := s.backend.deleteNode(targetUrl)
	if err != nil {
		return err
	}
	// no nodes left then remove server
	if len(s.backend.nodes) == 0 {
		delete(p.servers, serverName)
	}
	
	return nil
}
