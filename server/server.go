package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	statsGo "github.com/fukata/golang-stats-api-handler"
	"github.com/mercari/widebullet"
	"github.com/mercari/widebullet/jsonrpc"
	"github.com/mercari/widebullet/wlog"
)

var (
	HttpClient http.Client
)

func RegisterHandlers() {
	http.HandleFunc("/wbt", wideBulletHandler)
	statsGo.PrettyPrintEnabled()
	http.HandleFunc("/stat/go", statsGo.Handler)
}

func Run() {

	port := wbt.Config.Port

	HttpClient = http.Client{
		Timeout: time.Duration(wbt.Config.Timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: wbt.Config.MaxIdleConnsPerHost,
			DisableCompression:  wbt.Config.DisableCompression,
		},
	}

	// Listen TCP Port
	if _, err := strconv.Atoi(port); err == nil {
		errorLog(wlog.Debug, "listen port:%s", port)
		http.ListenAndServe(":"+port, nil)
	}

	// Listen UNIX Socket
	if strings.HasPrefix(port, "unix:/") {
		sockPath := port[5:]
		fi, err := os.Lstat(sockPath)
		if err == nil && (fi.Mode()&os.ModeSocket) == os.ModeSocket {
			err := os.Remove(sockPath)
			if err != nil {
				log.Fatal("failed to remove " + sockPath)
			}
		}
		l, err := net.Listen("unix", sockPath)
		if err != nil {
			log.Fatal("failed to listen: " + sockPath)
		}
		errorLog(wlog.Debug, "listen port:%s", port)
		http.Serve(l, nil)
	}

	errorLog(wlog.Error, "failed to listen port:%s", port)
}

func sendTextResponse(w http.ResponseWriter, result string, code int) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Server", wbt.ServerHeader())
	http.Error(w, result, code)
}

func sendJsonResponse(w http.ResponseWriter, result string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", wbt.ServerHeader())
	fmt.Fprint(w, result)
}

func wideBulletHandler(w http.ResponseWriter, r *http.Request) {
	var reqs []jsonrpc.Request

	if r.Method != "POST" {
		accessLog(r, &reqs)
		sendTextResponse(w, "method must be POST", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()
	if err := decoder.Decode(&reqs); err != nil {
		accessLog(r, &reqs)
		sendTextResponse(w, "request is malformed", http.StatusBadRequest)
		return
	}

	accessLog(r, &reqs)

	if err := jsonrpc.ValidateRequests(&reqs); err != nil {
		errorLog(wlog.Error, err.Error())
		sendTextResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resps, err := jsonRpc2Http(&reqs)
	if err != nil {
		errorLog(wlog.Error, err.Error())
		sendTextResponse(w, err.Error(), http.StatusBadGateway)
		return
	}

	bytes, err := json.Marshal(&resps)
	if err != nil {
		errorLog(wlog.Error, err.Error())
		sendTextResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJsonResponse(w, string(bytes))
}
