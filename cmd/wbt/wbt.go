package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mercari/widebullet"
	"github.com/mercari/widebullet/config"
	"github.com/mercari/widebullet/server"
	"github.com/mercari/widebullet/wlog"
)

func main() {
	versionPrinted := flag.Bool("v", false, "print widebullet version")
	port := flag.String("p", "", "listening port number or socket path")
	configPath := flag.String("c", "", "configuration file path")
	flag.Parse()

	if *versionPrinted {
		wbt.PrintVersion()
		return
	}

	conf, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// overwrite if port is specified by flags
	if *port != "" {
		conf.Port = *port
	}

	// set global configuration
	wbt.Config = conf
	wbt.AL = wlog.AccessLogger(conf.LogLevel)
	wbt.EL = wlog.ErrorLogger(conf.LogLevel)

	// Setup server
	mux := http.NewServeMux()
	server.RegisterHandlers(mux)
	server.SetupClient(&wbt.Config)

	srv := &http.Server{
		Handler: mux,
	}

	go func() {
		wbt.EL.Out(wlog.Debug, "Start running server")
		if err := server.Run(srv, &wbt.Config); err != nil {
			wbt.EL.Out(wlog.Error, "Failed to run server: %s", err)
		}
	}()

	// Watch SIGTERM signal and then gracefully shutdown
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

	wbt.EL.Out(wlog.Debug, "Start to shutdown server")
	timeout := time.Duration(conf.ShutdownTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		wbt.EL.Out(wlog.Error, "Failed to shutdown server: %s", err)
		return
	}

	wbt.EL.Out(wlog.Debug, "Successfully shutdown server")
}
