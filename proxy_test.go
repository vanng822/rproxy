package rproxy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProxy(t *testing.T) {
	assert := assert.New(t)
	p := NewProxy(nil)
	assert.True(p.conf.ApiEnable)
	assert.Equal(5555, p.conf.Port)
	assert.Nil(p.conf.Servers)
}
