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
