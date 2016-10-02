package server

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/fujiwara/fluent-agent-hydra/ltsv"
	"github.com/mercari/widebullet"
	"github.com/mercari/widebullet/jsonrpc"
	"github.com/mercari/widebullet/wlog"
)

func accessLog(r *http.Request, rr *[]jsonrpc.Request) {
	records := make(map[string]interface{})
	records["addr"] = r.RemoteAddr
	records["length"] = r.ContentLength
	if wbt.Config.LogLevel == "debug" {
		records["headers"] = r.Header
		records["body"] = *rr
	}
	buf := &bytes.Buffer{}
	encoder := ltsv.NewEncoder(buf)
	encoder.Encode(records)
	wbt.AL.Out(wlog.Info, strings.TrimRight(buf.String(), "\n"))
}

func errorLog(level wlog.LogLevel, msg string, args ...interface{}) {
	wbt.EL.Out(level, msg, args...)
}
