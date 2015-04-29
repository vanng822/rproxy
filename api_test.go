package rproxy

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAdminAPIWithDeleteBackendNode(t *testing.T) {
	assert := assert.New(t)
	p := NewProxy(nil)

	adminAPI := p.AdminAPI()
	ts := httptest.NewServer(adminAPI)
	defer ts.Close()

	client := &http.Client{}
	data := url.Values{}
	data.Add("serverName", "fakedomain.tld")
	data.Add("targetUrl", "127.0.0.1:8080")
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/_server/backend", ts.URL),
		bytes.NewBufferString(data.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.Nil(err)
	res, err := client.Do(req)
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
	content, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(content), "OK")

	// one 
	req, err = http.NewRequest("GET",
		fmt.Sprintf("%s/_server/backend?serverName=fakedomain.tld", ts.URL),
		nil)
	assert.Nil(err)
	res, err = client.Do(req)
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
	content, _ = ioutil.ReadAll(res.Body)
	assert.Contains(string(content), "127.0.0.1:8080")
	
	// all
	req, err = http.NewRequest("GET",
		fmt.Sprintf("%s/_server", ts.URL),
		nil)
	assert.Nil(err)
	res, err = client.Do(req)
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
	content, _ = ioutil.ReadAll(res.Body)
	assert.Contains(string(content), "127.0.0.1:8080")
	
	// Delete; data in url not body
	data = url.Values{}
	data.Add("serverName", "fakedomain.tld")
	data.Add("targetUrl", "127.0.0.1:8080")
	req, err = http.NewRequest("DELETE",
		fmt.Sprintf("%s/_server/backend?%s", ts.URL, data.Encode()),
		nil)
	assert.Nil(err)
	res, err = client.Do(req)
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
	content, _ = ioutil.ReadAll(res.Body)
	assert.Contains(string(content), "OK")
	
	// list empty
	req, err = http.NewRequest("GET",
		fmt.Sprintf("%s/_server", ts.URL),
		nil)
	assert.Nil(err)
	res, err = client.Do(req)
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
	content, _ = ioutil.ReadAll(res.Body)
	assert.NotContains(string(content), "127.0.0.1:8080")
}



func TestAdminAPIWithDeleteServerName(t *testing.T) {
	assert := assert.New(t)
	p := NewProxy(nil)

	adminAPI := p.AdminAPI()
	ts := httptest.NewServer(adminAPI)
	defer ts.Close()

	client := &http.Client{}
	data := url.Values{}
	data.Add("serverName", "fakedomain.tld")
	data.Add("targetUrl", "127.0.0.1:8080")
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/_server/backend", ts.URL),
		bytes.NewBufferString(data.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.Nil(err)
	res, err := client.Do(req)
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
	content, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(content), "OK")
	
	// all
	req, err = http.NewRequest("GET",
		fmt.Sprintf("%s/_server", ts.URL),
		nil)
	assert.Nil(err)
	res, err = client.Do(req)
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
	content, _ = ioutil.ReadAll(res.Body)
	assert.Contains(string(content), "127.0.0.1:8080")
	
	// Delete; data in url not body
	data = url.Values{}
	data.Add("serverName", "fakedomain.tld")
	req, err = http.NewRequest("DELETE",
		fmt.Sprintf("%s/_server?%s", ts.URL, data.Encode()),
		nil)
	assert.Nil(err)
	res, err = client.Do(req)
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
	content, _ = ioutil.ReadAll(res.Body)
	assert.Contains(string(content), "OK")
	
	// list empty
	req, err = http.NewRequest("GET",
		fmt.Sprintf("%s/_server", ts.URL),
		nil)
	assert.Nil(err)
	res, err = client.Do(req)
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
	content, _ = ioutil.ReadAll(res.Body)
	assert.NotContains(string(content), "127.0.0.1:8080")
}
