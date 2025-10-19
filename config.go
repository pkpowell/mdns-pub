package main

import (
	"errors"
	"io/fs"
	"os"

	"github.com/spf13/viper"
)

type HTTPServer struct {
	Address string
	Port    int
}
type Config struct {
	Servers    []Server
	HTTPServer HTTPServer
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add search paths to find the file
	viper.AddConfigPath("/etc/mdns-pub/")

	var fileNotFoundError viper.ConfigFileNotFoundError
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &fileNotFoundError) {
			Errorf("Config file not found %s", err)
			err = viper.WriteConfig()
			if err != nil {
				Errorf("viper.WriteConfig error %s", err)
				os.Exit(1)
			}
			// os.Exit(1)
			// Indicates an explicitly set config file is not found (such as with
			// using `viper.SetConfigFile`) or that no config file was found in
			// any search path (such as when using `viper.AddConfigPath`)
		} else {
			Errorf("Config file error %s", err)
			os.Exit(1)
			// Config file was found but another error was produced
		}
	}

	// viper.

	// if !exists(configFile) {
	// 	Info("Writing default config to", configFile)
	// 	err := viper.SafeWriteConfig()
	// 	if err != nil {
	// 		Error("viper.SafeWriteConfig error", err.Error())
	// 		a.terminate <- syscall.SIGINT
	// 	}
	// }

	viper.SetDefault("settings", &Config{
		Servers: []Server{},
		HTTPServer: HTTPServer{
			Address: "0.0.0.0",
			Port:    1122,
		},
	})

	err := viper.ReadInConfig()
	if err != nil {
		Errorf("fatal error config file: %s", err)
		return
	}
}

func exists(f string) bool {
	Debug("checking path", f)
	if _, err := os.Stat(f); errors.Is(err, fs.ErrNotExist) {
		Errorf("File not found", f)
		return false
	}
	// Debug("found", f)
	return true
}
