package rproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type BackendNode struct {
	targetUrl string
	server    *httputil.ReverseProxy
}

type Backend struct {
	// map on targetUrl
	nodes []*BackendNode
}

func (b *Backend) deleteNode(targetUrl string) error {
	found := -1
	for i, n := range b.nodes {
		if n.targetUrl == targetUrl {
			found = i
			break
		}
	}

	if found > 0 {
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
	if n == nil {
		n = &BackendNode{}
	}
	n.targetUrl = targetUrl
	n.server = httputil.NewSingleHostReverseProxy(target)
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
	// TODO round rubin
	for _, p := range b.nodes {
		return p
	}

	return nil
}

func (b *Backend) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	node := b.nextNode()
	if node != nil {
		// TODO next node if this node fails to serve
		node.server.ServeHTTP(rw, req)
	}
}
