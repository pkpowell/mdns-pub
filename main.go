package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grandcat/zeroconf"
)

type App struct {
	Servers  []*Server
	hostname string
}

type Server struct {
	Name        string
	Hostname    string
	Service     string
	IPAddress   string
	Port        int
	Extra       string
	MDNSService *zeroconf.Server
	// HCService   *mdns.Server
}

var servers = []*Server{
	{
		Name:      "WGS",
		IPAddress: "192.168.44.60",
		Port:      445,
		Service:   "_smb._tcp.",
		Hostname:  "wgs",
		Extra:     "WGS - Boltshauser Architekten AG",
	},
	{
		Name:      "Delta",
		IPAddress: "192.168.44.65",
		Port:      445,
		Service:   "_smb._tcp.",
		Hostname:  "delta",
		Extra:     "Delta - Boltshauser Architekten AG",
	},
	{
		Name:      "Bim Binz",
		IPAddress: "192.168.24.22",
		Port:      80,
		Service:   "_workstation._tcp.",
		Hostname:  "bim-binz",
		Extra:     "Bim Binz - Boltshauser Architekten AG",
	},
	{
		Name:      "adm",
		IPAddress: "192.168.44.61",
		Port:      445,
		Service:   "_smb._tcp.",
		Hostname:  "adm",
		Extra:     "adm / updates - Boltshauser Architekten AG",
	},
	{
		Name:      "adm",
		IPAddress: "192.168.44.61",
		Port:      1212,
		Service:   "_ztui._tcp.",
		Hostname:  "adm",
		Extra:     "adm / ztui - Boltshauser Architekten AG",
	},
	{
		Name:      "zt-binz",
		IPAddress: "192.168.24.41",
		Port:      1212,
		Service:   "_ztui._tcp.",
		Hostname:  "zt-binz",
		Extra:     "adm / ztui - Boltshauser Architekten AG",
	},
	{
		Name:      "adm",
		IPAddress: "192.168.44.61",
		Port:      22,
		Service:   "_ssh._tcp.",
		Hostname:  "adm",
		Extra:     "adm / ssh - Boltshauser Architekten AG",
	},
	{
		Name:      "adm",
		IPAddress: "192.168.44.61",
		Port:      5900,
		Service:   "_rfb._tcp.",
		Hostname:  "adm",
		Extra:     "adm / screen sharing - Boltshauser Architekten AG",
	},
}

var (
	terminate = make(chan os.Signal, 1)
	update    = make(chan bool, 1)
)

func main() {
	var err error
	var a = &App{
		Servers: servers,
	}

	signal.Notify(terminate,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	initLogging()
	initConfig()

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

			for _, server := range servers {
				if server.MDNSService != nil {
					Infof("Canceling %s for %s", server.Service, server.Name)
					server.MDNSService.Shutdown()
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

	for _, server := range a.Servers {

		server.MDNSService, err = zeroconf.RegisterProxy(
			server.Name,
			server.Service,
			"local.",
			server.Port,
			server.Hostname,
			[]string{server.IPAddress},
			[]string{server.Extra, "published by " + a.hostname},
			[]net.Interface{*iface},
		)

		if err != nil {
			Errorf("zeroconf.Register error %s", err)
			return
		}

		Infof("Publishing %s for %s on %s", server.Service, server.Name, iface.Name)
	}
}
