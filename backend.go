package rproxy

import (
	"net/http"
	"net/http/httputil"
)

type Backend struct {
	// map on targetUrl
	nodes map[string]*httputil.ReverseProxy
}

func (b *Backend) getNode() *httputil.ReverseProxy {
	// Pickup one backend proxy
	// TODO round rubin
	for _, p := range b.nodes {
		return p
	}

	return nil
}

func (b *Backend) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	node := b.getNode()
	if node != nil {
		// TODO next node if this node fails to serve
		node.ServeHTTP(rw, req)
	}
}