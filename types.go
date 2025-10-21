package main

import "github.com/grandcat/zeroconf"

type App struct {
	// servers  []*Server
	hostname string
	Config   *Config `json:"config,omitempty"`
}

type Server struct {
	Name      string           `json:"name,omitempty"`
	Hostname  string           `json:"hostname,omitempty"`
	Service   string           `json:"service,omitempty"`
	IPAddress string           `json:"ipAddress,omitempty"`
	Port      int              `json:"port,omitempty"`
	ExtraData string           `json:"extraData,omitempty"`
	zeroconf  *zeroconf.Server `json:"-"`
	// HCService   *mdns.Server
}

type Conf struct {
	Config *Config `json:"config"`
}

type HTTPServer struct {
	Address string `json:"address,omitempty"`
	Port    int    `json:"port,omitempty"`
}

type Config struct {
	Servers    []*Server  `json:"servers,omitempty"`
	HTTPServer HTTPServer `json:"httpServer,omitempty"`
}
