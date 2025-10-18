package main

import (
	"github.com/pkpowell/logging"
)

var (
	Error, Warn, Info, Debug     logging.Display
	Infop                        logging.Display
	Errorf, Warnf, Infof, Debugf logging.Displayf
	Errorw                       logging.Create
)

func initLogging() {
	logging.Init(boolPointer(false), boolPointer(false), boolPointer(true))
	// logging.Init(verbose, jsonLogs, colour)

	Error = logging.Error
	Warn = logging.Warn
	Info = logging.Info
	Infop = logging.Infop
	Debug = logging.Debug

	Errorf = logging.Errorf
	Warnf = logging.Warnf
	Infof = logging.Infof
	Debugf = logging.Debugf

	Errorw = logging.Errorw
}

func boolPointer(b bool) *bool {
	return &b
}
