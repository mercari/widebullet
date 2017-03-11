package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mercari/widebullet"
	"github.com/mercari/widebullet/config"
	"github.com/mercari/widebullet/jsonrpc"
	"github.com/stretchr/testify/assert"
)

func TestBuildRequestURI(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("http://127.0.0.1:30001/resource/get?id=1",
		buildRequestURI("127.0.0.1:30001", "/resource/get", "?id=1"))
	assert.Equal("http://127.0.0.1:30001/resource/get?id=1",
		buildRequestURI("http://127.0.0.1:30001", "/resource/get", "?id=1"))
}

func TestBuildHttpRequest(t *testing.T) {
	assert := assert.New(t)

	payload := `[
    {"jsonrpc": "2.0", "ep": "ep-1",                        "method": "/user/get",    "params": { "user_id": 1 },                   "id": "1"},
    {"jsonrpc": "2.0", "ep": "ep-1", "http_method": "GET",  "method": "/item/get",    "params": { "item_id": 2 },                   "id": "2"},
    {"jsonrpc": "2.0", "ep": "ep-2", "http_method": "POST", "method": "/item/update", "params": { "item_id": 2, "desc": "update" }, "id": "3"}
]
`

	var (
		reqjs []jsonrpc.Request
		reqhs []http.Request
		buf   bytes.Buffer
	)
	decoder := json.NewDecoder(strings.NewReader(payload))
	decoder.UseNumber()
	err := decoder.Decode(&reqjs)
	assert.Nil(err)

	wbt.Config, err = config.Load("../config/example.toml")
	assert.Nil(err)

	headers := make(http.Header)
	headers.Add("X-Forwarded-For", "127.0.0.1")
	headers.Add("X-Auth-Token", "Bearer test_token")
	headers.Add("X-Auth-Token2", "Bearer another_testtoken")
	for _, reqj := range reqjs {
		reqh, err := buildHttpRequest(&reqj, &headers)
		assert.Nil(err)
		assert.Equal(wbt.ServerHeader(), reqh.Header.Get("User-Agent"))
		assert.Equal("127.0.0.1", reqh.Header.Get("X-Forwarded-For"))
		reqhs = append(reqhs, *reqh)
	}

	assert.Equal("", reqhs[0].Header.Get("Content-Type"))
	assert.Equal("ep1.example.com", reqhs[0].Host)
	assert.Equal("127.0.0.1:30001", reqhs[0].URL.Host)
	assert.Equal("/user/get", reqhs[0].URL.Path)
	assert.Equal("user_id=1", reqhs[0].URL.RawQuery)
	assert.Equal("Bearer test_token", reqhs[0].Header.Get("Authorization"))

	assert.Equal("", reqhs[1].Header.Get("Content-Type"))
	assert.Equal("ep1.example.com", reqhs[1].Host)
	assert.Equal("127.0.0.1:30001", reqhs[1].URL.Host)
	assert.Equal("/item/get", reqhs[1].URL.Path)
	assert.Equal("item_id=2", reqhs[1].URL.RawQuery)
	assert.Equal("Bearer test_token", reqhs[1].Header.Get("Authorization"))

	buf.ReadFrom(reqhs[2].Body)
	assert.Equal("application/x-www-form-urlencoded", reqhs[2].Header.Get("Content-Type"))
	assert.Equal("ep2.example.com", reqhs[2].Host)
	assert.Equal("127.0.0.1:30002", reqhs[2].URL.Host)
	assert.Equal("/item/update", reqhs[2].URL.Path)
	assert.Equal("desc=update&item_id=2", buf.String())
	assert.Equal("Bearer another_testtoken", reqhs[2].Header.Get("Authorization"))
}
