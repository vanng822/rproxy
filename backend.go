package rproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func CreateReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

type BackendNode struct {
	targetUrl string
	server    *httputil.ReverseProxy
}

type Backend struct {
	nodes []*BackendNode
	next  int
}

func NewBackend() *Backend {
	b := &Backend{}
	b.nodes = make([]*BackendNode, 0)
	return b
}

func (b *Backend) deleteNode(targetUrl string) error {
	found := -1
	for i, n := range b.nodes {
		if n.targetUrl == targetUrl {
			found = i
			break
		}
	}

	if found > -1 {
		if found == len(b.nodes)-1 {
			b.nodes = append(b.nodes[:found])
		} else {
			b.nodes = append(b.nodes[:found], b.nodes[found+1:]...)
		}
		return nil
	}

	return fmt.Errorf("Could not found node %s", targetUrl)
}

func (b *Backend) addNode(targetUrl string) error {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return err
	}
	// find if there is a node for this targetUrl
	n := b.getNode(targetUrl)
	if n != nil {
		return nil
	}
	n = &BackendNode{}
	n.targetUrl = targetUrl
	n.server = CreateReverseProxy(target)
	b.nodes = append(b.nodes, n)
	return nil
}

func (b *Backend) getNode(targetUrl string) *BackendNode {
	for _, n := range b.nodes {
		if n.targetUrl == targetUrl {
			return n
		}
	}
	return nil
}

// nextNode returns a node candidate for serving the request
func (b *Backend) nextNode() *BackendNode {
	// Pickup one backend proxy
	if len(b.nodes) == 0 {
		return nil
	}
	p := b.nodes[b.next]
	b.next++
	if b.next >= len(b.nodes) {
		b.next = 0
	}

	return p
}

func (b *Backend) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	node := b.nextNode()
	if node != nil {
		// TODO next node if this node fails to serve
		node.server.ServeHTTP(rw, req)
	} else {
		http.Error(rw, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	}

}
