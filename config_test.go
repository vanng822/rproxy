package rproxy

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	assert := assert.New(t)
	config := strings.NewReader(`{
		"Host": "127.0.0.1",
		"Port": 8080,
		"Servers": [
			{"ServerName": "dev.com", "TargetUrl": "127.0.0.1:8090"},
			{"ServerName": "dev.com", "TargetUrl": "127.0.0.1:8091"}
		]
	}`)
	conf := ParseConfig(config)
	assert.Equal("127.0.0.1", conf.Host)
	assert.Equal(8080, conf.Port)
	assert.Equal(2, len(conf.Servers))
	assert.Equal("dev.com", conf.Servers[0].ServerName)
	assert.Equal("127.0.0.1:8090", conf.Servers[0].TargetUrl)
}
