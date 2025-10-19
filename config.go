package main

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

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

func (a *App) initConfig() {
	var err error
	// var fileNotFoundError viper.ConfigFileNotFoundError
	var fileExistsError viper.ConfigFileAlreadyExistsError

	var configPath = "/etc/mdns-pub/"
	var configFileName = "config"
	var configFilenType = "json"

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFilenType)

	viper.AddConfigPath(configPath)

	// viper.Set("config", &Config{
	// 	Servers: []*Server{
	// 		{
	// 			Name:      "server-name",
	// 			Hostname:  "server-name.local",
	// 			Port:      1234,
	// 			Service:   "_workstation._tcp.",
	// 			IPAddress: "192.168.22.1",
	// 			ExtraData: "",
	// 		},
	// 	},
	// 	HTTPServer: HTTPServer{
	// 		Address: "0.0.0.0",
	// 		Port:    1122,
	// 	},
	// })

	Info("Writing default config to", viper.ConfigFileUsed())
	err = viper.SafeWriteConfig()
	if err != nil {
		if !errors.As(err, &fileExistsError) {
			Error("viper.SafeWriteConfig error", err.Error())
			os.Exit(1)
		}
	}

	err = viper.ReadInConfig()
	if err != nil {
		Errorf("fatal error config file: %s", err)
		return
	}

	var conf Conf

	err = viper.Unmarshal(&conf)
	if err != nil {
		Error("viper.Unmarshal", err.Error())
		os.Exit(1)
	}

	Infof("Config loaded %v", conf)
	a.Config = conf.Config
}

// func exists(f string) bool {
// 	Debug("checking path", f)
// 	if _, err := os.Stat(f); errors.Is(err, fs.ErrNotExist) {
// 		Errorf("File not found %s", f)
// 		return false
// 	}
// 	// Debug("found", f)
// 	return true
// }
