package rproxy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackendAddNode(t *testing.T) {
	assert := assert.New(t)
	b := NewBackend()
	b.addNode("http://127.0.0.1:8080")
	b.addNode("http://127.0.0.1:8080")
	assert.Len(b.nodes, 1)
}

func TestBackendDeleteNode(t *testing.T) {
	assert := assert.New(t)
	b := NewBackend()
	b.addNode("http://127.0.0.1:8080")
	assert.Len(b.nodes, 1)
	b.deleteNode("http://127.0.0.1:8080")
	assert.Len(b.nodes, 0)
}

func TestBackendNextNode(t *testing.T) {
	assert := assert.New(t)
	nodes := []string{"http://127.0.0.1:8080", "http://127.0.0.1:8081", "http://127.0.0.1:8082"}
	b := NewBackend()
	
	for _, n := range nodes {
		b.addNode(n)
	}
	assert.Len(b.nodes, 3)
	
	for i := 0; i < 3; i++ {
		for _, n := range nodes {
			node := b.nextNode()
			assert.Equal(node.targetUrl, n)
		}
	}
}
