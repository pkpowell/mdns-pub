package main

import (
	"net"
	"os"

	"github.com/grandcat/zeroconf"
)

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
				if server.zeroconf != nil {
					Infof("Canceling %s for %s", server.Service, server.Name)
					server.zeroconf.Shutdown()
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

		server.zeroconf, err = zeroconf.RegisterProxy(
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
