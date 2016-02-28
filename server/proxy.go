package server

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/mercari/widebullet/jsonrpc"
	"github.com/mercari/widebullet/wlog"
)

func jsonRpc2Http(reqs *[]jsonrpc.Request) ([]jsonrpc.Response, error) {
	wg := new(sync.WaitGroup)
	resps := make([]jsonrpc.Response, len(*reqs))
	// send requests to endpoint conccurrently
	for i, reqj := range *reqs {
		wg.Add(1)
		go func(i int, reqj jsonrpc.Request) {
			defer wg.Done()
			reqh, err := buildHttpRequest(&reqj)
			if err != nil {
				resps[i] = buildJsonRpcErrorResponse(jsonrpc.InternalError, err.Error(), reqj.ID, 0)
				errorLog(wlog.Error, err.Error())
				return
			}
			start := time.Now()
			resp, err := HttpClient.Do(reqh)
			end := time.Now()
			ptime := (end.Sub(start)).Seconds()
			if err != nil {
				resps[i] = buildJsonRpcErrorResponse(jsonrpc.InternalError, err.Error(), reqj.ID, ptime)
				errorLog(wlog.Error, err.Error())
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				resps[i] = buildHttpError2JsonRpcErrorResponse(resp, reqj.ID, ptime)
				errorLog(wlog.Error, "%#v is failed: %s", reqj, resp.Status)
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				resps[i] = buildJsonRpcErrorResponse(jsonrpc.InternalError, err.Error(), reqj.ID, ptime)
				errorLog(wlog.Error, err.Error())
				return
			}
			resps[i] = buildJsonRpcResponse(string(body), reqj.ID, ptime)
			return
		}(i, reqj)
	}

	wg.Wait()

	return resps, nil
}
