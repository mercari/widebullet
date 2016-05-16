package server

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/mercari/widebullet/jsonrpc"
	"github.com/mercari/widebullet/wlog"
)

func sendHttpRequest(wg *sync.WaitGroup, reqj jsonrpc.Request, forwardHeaders *http.Header, respj *jsonrpc.Response) {
	defer wg.Done()
	reqh, err := buildHttpRequest(&reqj, forwardHeaders)
	if err != nil {
		*respj = buildJsonRpcErrorResponse(jsonrpc.InternalError, err.Error(), reqj.ID, 0)
		errorLog(wlog.Error, err.Error())
		return
	}
	start := time.Now()
	resp, err := HttpClient.Do(reqh)
	end := time.Now()
	ptime := (end.Sub(start)).Seconds()
	if err != nil {
		*respj = buildJsonRpcErrorResponse(jsonrpc.InternalError, err.Error(), reqj.ID, ptime)
		errorLog(wlog.Error, err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		*respj = buildHttpError2JsonRpcErrorResponse(resp, reqj.ID, ptime)
		errorLog(wlog.Error, "%#v is failed: %s", reqj, resp.Status)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		*respj = buildJsonRpcErrorResponse(jsonrpc.InternalError, err.Error(), reqj.ID, ptime)
		errorLog(wlog.Error, err.Error())
		return
	}
	*respj = buildJsonRpcResponse(string(body), reqj.ID, ptime)
}

func jsonRpc2Http(reqs *[]jsonrpc.Request, forwardHeaders *http.Header) ([]jsonrpc.Response, error) {
	wg := new(sync.WaitGroup)
	resps := make([]jsonrpc.Response, len(*reqs))
	// send requests to endpoint conccurrently
	for i, reqj := range *reqs {
		wg.Add(1)
		go sendHttpRequest(wg, reqj, forwardHeaders, &resps[i])
	}

	wg.Wait()

	return resps, nil
}
