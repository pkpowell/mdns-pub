package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grandcat/zeroconf"
)

type App struct {
	// servers  []*Server
	hostname string
	Config   *Config `json:"config,omitempty"`
}

type Server struct {
	Name        string           `json:"name,omitempty"`
	Hostname    string           `json:"hostname,omitempty"`
	Service     string           `json:"service,omitempty"`
	IPAddress   string           `json:"ipAddress,omitempty"`
	Port        int              `json:"port,omitempty"`
	ExtraData   string           `json:"extraData,omitempty"`
	mdnsService *zeroconf.Server `json:"-"`
	// HCService   *mdns.Server
}

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

func (a *App) initMDNS() {
	var err error
	var iface net.Interface

	ifaces, err := net.Interfaces()
	if err != nil {
		Errorf("net.Interfaces error %s", err.Error())
		return
	}

	for _, iface = range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			Infof("found loopback interface %s", iface.Name)
			go a.publish(&iface)
			break
		}
	}

	for {
		select {
		case signal := <-terminate:
			Info("Received terminate signal", signal.String())
			Infof("cleaning up...")

			for _, server := range a.Config.Servers {
				if server.mdnsService != nil {
					Infof("Canceling %s for %s", server.Service, server.Name)
					server.mdnsService.Shutdown()
				}
			}
			os.Exit(0)

		case <-update:
			// do updates...
		}
	}
}

func (a *App) publish(iface *net.Interface) {
	var err error

	if len(a.Config.Servers) == 0 {
		Warnf("no servers configured")
		return
	}

	for _, server := range a.Config.Servers {

		server.mdnsService, err = zeroconf.RegisterProxy(
			server.Name,
			server.Service,
			"local.",
			server.Port,
			server.Hostname,
			[]string{server.IPAddress},
			[]string{server.ExtraData, "published by " + a.hostname},
			[]net.Interface{*iface},
		)

		if err != nil {
			Errorf("zeroconf.Register error %s", err)
			return
		}

		Infof("Publishing %s for %s on %s", server.Service, server.Name, iface.Name)
	}
}
