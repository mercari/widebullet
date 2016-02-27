package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadExampleToml(t *testing.T) {
	assert := assert.New(t)

	c, err := Load("./example.toml")
	assert.Nil(err)

	assert.Equal("29300", c.Port)
	assert.Equal("error", c.LogLevel)
	assert.Equal(5, c.Timeout)
	assert.Equal(100, c.MaxIdleConnsPerHost)
	assert.Equal(false, c.DisableCompression)

	eps := c.Endpoints
	assert.Equal(2, len(eps))
}

func TestFindEp(t *testing.T) {
	assert := assert.New(t)

	c, err := Load("./example.toml")
	assert.Nil(err)

	ep, err := FindEp(c, "ep-1")
	assert.Nil(err)
	assert.Equal("ep-1", ep.Name)
	assert.Equal("127.0.0.1:30001", ep.Ep)
	assert.Equal("Host", ep.ProxySetHeaders[0][0])
	assert.Equal("ep1.example.com", ep.ProxySetHeaders[0][1])

	ep, err = FindEp(c, "ep-2")
	assert.Nil(err)
	assert.Equal("ep-2", ep.Name)
	assert.Equal("127.0.0.1:30002", ep.Ep)
	assert.Equal("Host", ep.ProxySetHeaders[0][0])
	assert.Equal("ep2.example.com", ep.ProxySetHeaders[0][1])

	_, err = FindEp(c, "ep-3")
	assert.NotNil(err)
}
