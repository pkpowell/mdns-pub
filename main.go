package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grandcat/zeroconf"
)

type App struct {
	Servers []Server
}

type Server struct {
	Name        string
	Hostname    string
	Service     string
	IPAddress   string
	Port        int
	Extra       string
	MDNSService *zeroconf.Server
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

	signal.Notify(terminate,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	initLogging()

	initMDNS()

}

func initMDNS() {
	var err error
	var i net.Interface
	ifaces, err := net.Interfaces()
	if err != nil {
		Errorf("net.Interfaces error %s", err.Error())
		return
	}
cont:
	for _, i = range ifaces {
		Infof("iface %#v", i)

		switch i.Name {
		case "lo", "lo0":
			Infof("found loopback %s", i.Name)
			break cont
		default:
			Warn("No localhost interface found")
			// terminate <- syscall.SIGQUIT
			return
		}
	}

	for _, server := range servers {

		server.MDNSService, err = zeroconf.RegisterProxy(
			server.Name,
			server.Service,
			"local.",
			server.Port,
			server.Hostname,
			[]string{server.IPAddress},
			[]string{server.Extra},
			[]net.Interface{i},
		)

		if err != nil {
			Errorf("zeroconf.Register error %s", err)

			return
		}

		Infof("Publishing %s for %s", server.Service, server.Name)

	}

	for {
		select {
		case signal := <-terminate:
			Info("Received terminate signal", signal.String())
			Infof("cleaning up...")

			for _, server := range servers {
				Infof("Canceling %s for %s", server.Service, server.Name)
				server.MDNSService.Shutdown()
			}
			os.Exit(0)

		case <-update:
			// do updates...
		}
	}

	// ifa, err := net.InterfaceByName("lo0")
	// if err != nil {
	// 	Errorf("net.InterfaceByName error %s", err)
	// 	return
	// }

}
