package main

import (
	"flag"
	"log"

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

	server.RegisterHandlers()
	server.Run()
}
