package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mercari/widebullet"
	"github.com/mercari/widebullet/config"
	"github.com/mercari/widebullet/jsonrpc"
)

func buildRequestURI(ep, method, qs string) string {
	if strings.HasPrefix(ep, "http://") {
		return fmt.Sprintf("%s%s%s", ep, method, qs)
	}
	if strings.HasPrefix(ep, "https://") {
		return fmt.Sprintf("%s%s%s", ep, method, qs)
	}
	return fmt.Sprintf("http://%s%s%s", ep, method, qs)
}

func buildURLEncodedString(params jsonrpc.RequestParams, method string) (string, error) {
	values := url.Values{}
	for k, v := range params {
		switch v.(type) {
		case string:
			values.Set(k, v.(string))
		case json.Number:
			values.Set(k, fmt.Sprintf("%s", v))
		default:
			return "", fmt.Errorf("wbt supports only string and number")
		}
	}

	if method == "POST" {
		return values.Encode(), nil
	}

	return fmt.Sprintf("?%s", values.Encode()), nil
}

func buildJsonRpcResponse(body, id string, time float64) jsonrpc.Response {
	return jsonrpc.Response{
		Version: jsonrpc.Version,
		Result:  body,
		ID:      id,
		Time:    time,
	}
}

func buildJsonRpcErrorResponse(code int, msg, id string, time float64) jsonrpc.Response {
	jsonRpcError := &jsonrpc.Error{
		Code:    code,
		Message: msg,
	}

	return jsonrpc.Response{
		Version: jsonrpc.Version,
		Error:   jsonRpcError,
		ID:      id,
		Time:    time,
	}
}

func buildHttpError2JsonRpcErrorResponse(status_code int, err_msg string, id string, time float64) jsonrpc.Response {

	r := &struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	} {
		status_code,
		err_msg,
	}

	jsonrpc_err_msg, err := json.Marshal(r)
	if err != nil {
		return buildJsonRpcErrorResponse(jsonrpc.InternalError, err.Error(), id, time)
	}

	switch status_code {
	case http.StatusNotFound:
		return buildJsonRpcErrorResponse(jsonrpc.MethodNotFoundError, string(jsonrpc_err_msg), id, time)
	}
	return buildJsonRpcErrorResponse(jsonrpc.InternalError, string(jsonrpc_err_msg), id, time)
}

func buildHttpRequest(reqj *jsonrpc.Request, forwardHeaders *http.Header) (*http.Request, error) {
	var reqh *http.Request

	ep, err := config.FindEp(wbt.Config, reqj.Ep)
	if err != nil {
		return reqh, err
	}

	es, err := buildURLEncodedString(reqj.Params, reqj.HttpMethod)
	if err != nil {
		return reqh, err
	}

	switch reqj.HttpMethod {
	case "POST":
		uri := buildRequestURI(ep.Ep, reqj.Method, "")
		reqh, err = http.NewRequest("POST", uri, strings.NewReader(es))
		if err != nil {
			return reqh, err
		}
		reqh.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	default:
		uri := buildRequestURI(ep.Ep, reqj.Method, es)
		reqh, err = http.NewRequest("GET", uri, nil)
		if err != nil {
			return reqh, err
		}
	}

	ua := forwardHeaders.Get("User-Agent")
	if ua == "" {
		reqh.Header.Set("User-Agent", wbt.ServerHeader())
	} else {
		reqh.Header.Set("User-Agent", ua)
	}

	xForwardedFor := forwardHeaders.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		reqh.Header.Set("X-Forwarded-For", xForwardedFor)
	}

	for _, headers := range ep.ProxySetHeaders {
		if len(headers) < 2 {
			continue
		}
		key := headers[0]
		value := strings.Join(headers[1:], ",")
		if key == "Host" {
			reqh.Host = value
		} else {
			reqh.Header.Set(key, value)
		}
	}

	return reqh, nil
}
