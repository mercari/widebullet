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
	assert.Equal(30, c.IdleConnTimeout)
	assert.Equal(60, c.ProxyReadTimeout)
	assert.Equal(15, c.ShutdownTimeout)

	eps := c.Endpoints
	assert.Equal(2, len(eps))
}

func TestLoadGlobalConfig(t *testing.T) {
	assert := assert.New(t)

	configStr := `
Port = "12345"
LogLevel = "debug"
Timeout = 10
MaxIdleConnsPerHost = 1000
DisableCompression = true
IdleConnTimeout = 90
ProxyReadTimeout = 120
`

	c, err := LoadBytes([]byte(configStr))
	assert.Nil(err)

	assert.Equal("12345", c.Port)
	assert.Equal("debug", c.LogLevel)
	assert.Equal(10, c.Timeout)
	assert.Equal(1000, c.MaxIdleConnsPerHost)
	assert.Equal(true, c.DisableCompression)
	assert.Equal(90, c.IdleConnTimeout)
	assert.Equal(120, c.ProxyReadTimeout)
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
	assert.Equal("Authorization", ep.ProxyPassHeaders[0][0])
	assert.Equal("X-Auth-Token", ep.ProxyPassHeaders[0][1])

	ep, err = FindEp(c, "ep-2")
	assert.Nil(err)
	assert.Equal("ep-2", ep.Name)
	assert.Equal("http://127.0.0.1:30002", ep.Ep)
	assert.Equal("Host", ep.ProxySetHeaders[0][0])
	assert.Equal("ep2.example.com", ep.ProxySetHeaders[0][1])
	assert.Equal("Authorization", ep.ProxyPassHeaders[0][0])
	assert.Equal("X-Auth-Token2", ep.ProxyPassHeaders[0][1])

	_, err = FindEp(c, "ep-3")
	assert.NotNil(err)
}
