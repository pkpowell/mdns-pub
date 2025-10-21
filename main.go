package main

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	terminate = make(chan os.Signal, 1)
	update    = make(chan bool, 1)
)

func main() {
	initLogging()
	var err error
	var a = &App{}
	a.initConfig()

	signal.Notify(terminate,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	a.hostname, err = os.Hostname()
	if err != nil {
		Infof("os Hostname error %s", err)
		return
	}

	a.initMDNS()
}
