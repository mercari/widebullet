package jsonrpc

import (
	"fmt"
)

const Version = "2.0"

const (
	ParseError          = -32700
	InvalidRequestError = -32600
	MethodNotFoundError = -32601
	InvalidParamsError  = -32602
	InternalError       = -32603
)

type RequestParams map[string]interface{}

type Request struct {
	Version    string        `json:"jsonrpc"`
	Method     string        `json:"method"`
	HttpMethod string        `json:"http_method"`
	Params     RequestParams `json:"params,omitempty"`
	ID         string        `json:"id"`
	// extention
	Ep string `json:"ep"`
}

type Response struct {
	Version string  `json:"jsonrpc"`
	Result  string  `json:"result,omitempty"`
	Error   *Error  `json:"error,omitempty"`
	ID      string  `json:"id"`
	Time    float64 `json:"time,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	//Data    interface{} `json:"data"`
}

func ValidateRequests(reqs *[]Request) error {
	idMap := make(map[string]bool, len(*reqs))
	for _, r := range *reqs {
		if _, ok := idMap[r.ID]; ok {
			return fmt.Errorf("ID:%s is duplicated.", r.ID)
		}
		if err := validateRequest(&r); err != nil {
			return err
		}
		idMap[r.ID] = true
	}
	return nil
}

func validateRequest(r *Request) error {
	if r.Version != Version {
		return fmt.Errorf("malformed JSON-RPC version: %s", r.Version)
	}
	if r.Method == "" {
		return fmt.Errorf("empty method")
	}
	// empty method is treated as GET.
	if r.HttpMethod != "" && r.HttpMethod != "GET" && r.HttpMethod != "POST" {
		return fmt.Errorf("malformed HTTP method: %s", r.HttpMethod)
	}
	if r.ID == "" {
		return fmt.Errorf("empty id")
	}
	return nil
}
