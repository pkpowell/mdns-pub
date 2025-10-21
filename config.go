package main

import (
	"errors"
	"os"
	"path"

	"github.com/spf13/viper"
)

var defaultConf = &Config{
	Servers: []*Server{
		{
			Name:      "Example Server",
			Hostname:  "example-hostname.local",
			Port:      1234,
			Service:   "_workstation._tcp.",
			IPAddress: "192.168.22.1",
			ExtraData: "Example Server",
		},
	},
	HTTPServer: HTTPServer{
		Address: "0.0.0.0",
		Port:    1122,
	},
}

func (a *App) initConfig() {
	var err error
	var fileExistsError viper.ConfigFileAlreadyExistsError

	var configPath = "/etc/mdns-pub/"
	var configFileName = "config"
	var configFileType = "json"
	var configFile = path.Join(configPath, configFileName+"."+configFileType)

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)

	viper.AddConfigPath(configPath)

	if !exists(configFile) {
		viper.Set("config", defaultConf)
		Info("Writing default config to", configFile)
		err = viper.SafeWriteConfig()
		if err != nil {
			if !errors.As(err, &fileExistsError) {
				Error("viper.SafeWriteConfig error", err.Error())
				os.Exit(1)
			}
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
