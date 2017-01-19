package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	statsGo "github.com/fukata/golang-stats-api-handler"
	"github.com/lestrrat/go-server-starter/listener"
	"github.com/mercari/widebullet"
	"github.com/mercari/widebullet/config"
	"github.com/mercari/widebullet/jsonrpc"
	"github.com/mercari/widebullet/wlog"
)

var (
	HttpClient http.Client
)

// RegisterHandlers sets handler to serve.
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/wbt", wideBulletHandler)

	statsGo.PrettyPrintEnabled()
	mux.HandleFunc("/stat/go", statsGo.Handler)
}

// SetupClient setups http.Client (which is globally used in this package)
// with given config.
func SetupClient(config *config.Config) {
	HttpClient = http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConnsPerHost:   wbt.Config.MaxIdleConnsPerHost,
			DisableCompression:    wbt.Config.DisableCompression,
			IdleConnTimeout:       time.Duration(config.IdleConnTimeout) * time.Second,
			ResponseHeaderTimeout: time.Duration(config.ProxyReadTimeout) * time.Second,
		},
	}
}

// Run starts the given server. By default, it tries to accept
// requests from `go-server-starter`. If not then check config.Port
// value. If nothing to listen, then returns error.
func Run(server *http.Server, config *config.Config) error {

	// If ServerStarterEnv is found, then use listerner from
	// `go-server-starter`. Even if it fails, not terminate here
	// but use the given config.Port
	if v := os.Getenv(listener.ServerStarterEnvVarName); len(v) != 0 {
		listeners, err := listener.ListenAll()
		if err != nil {
			errorLog(wlog.Error, "Failed to get listeners from go-server-starter: %s", err)
		} else {
			if len(listeners) == 0 {
				errorLog(wlog.Error, "No listener to listen is found")
			} else {
				errorLog(wlog.Debug, "Start accepting request: %s", listeners[0].Addr())
				return server.Serve(listeners[0])
			}
		}
	}

	port := config.Port
	if len(port) == 0 {
		return fmt.Errorf("no port to listen")
	}

	// Listen TCP Port
	if _, err := strconv.Atoi(port); err == nil {
		errorLog(wlog.Debug, "Start listening: %s", port)
		server.Addr = ":" + port
		return server.ListenAndServe()
	}

	// Listen UNIX Socket
	if strings.HasPrefix(port, "unix:/") {
		sockPath := port[5:]
		fi, err := os.Lstat(sockPath)
		if err == nil && (fi.Mode()&os.ModeSocket) == os.ModeSocket {
			err := os.Remove(sockPath)
			if err != nil {
				return fmt.Errorf("failed to remove socket: %s", sockPath)
			}
		}

		l, err := net.Listen("unix", sockPath)
		if err != nil {
			return fmt.Errorf("failed to listen socket %q: %s", sockPath, err)
		}

		errorLog(wlog.Debug, "Start accepting request: %s", l.Addr())
		return server.Serve(l)
	}

	return fmt.Errorf("failed to listen port: %s", port)
}

func sendTextResponse(w http.ResponseWriter, result string, code int) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Server", wbt.ServerHeader())
	w.WriteHeader(code)
	fmt.Fprint(w, result)
}

func sendJsonResponse(w http.ResponseWriter, result string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", wbt.ServerHeader())
	fmt.Fprint(w, result)
}

func wideBulletHandler(w http.ResponseWriter, r *http.Request) {
	var reqs []jsonrpc.Request

	stime := time.Now()

	if r.Method != "POST" {
		accessLog(r, &reqs, stime, http.StatusBadRequest)
		sendTextResponse(w, "method must be POST", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()
	if err := decoder.Decode(&reqs); err != nil {
		accessLog(r, &reqs, stime, http.StatusBadRequest)
		sendTextResponse(w, "request is malformed", http.StatusBadRequest)
		return
	}

	if err := jsonrpc.ValidateRequests(&reqs); err != nil {
		accessLog(r, &reqs, stime, http.StatusBadRequest)
		errorLog(wlog.Error, err.Error())
		sendTextResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resps, err := jsonRpc2Http(&reqs, &r.Header)
	if err != nil {
		accessLog(r, &reqs, stime, http.StatusBadGateway)
		errorLog(wlog.Error, err.Error())
		sendTextResponse(w, err.Error(), http.StatusBadGateway)
		return
	}

	bytes, err := json.Marshal(&resps)
	if err != nil {
		accessLog(r, &reqs, stime, http.StatusInternalServerError)
		errorLog(wlog.Error, err.Error())
		sendTextResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJsonResponse(w, string(bytes))

	accessLog(r, &reqs, stime, http.StatusOK)
}
