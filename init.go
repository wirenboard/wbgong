package wbgo

import (
	"errors"
	"log"
	"plugin"
)

var errResolve = errors.New("Could not resolve symbol")

var plug *plugin.Plugin
var (
	// Error logger
	Error *log.Logger
	// Warn logger
	Warn *log.Logger
	// Info logger
	Info *log.Logger
	// Debug logger
	Debug *log.Logger
)

// Init tries to load shared library
func Init(path string) (err error) {
	plug, err = plugin.Open(path)
	if err != nil {
		return err
	}
	Error, err = initLogger("Error")
	if err != nil {
		return err
	}
	Warn, err = initLogger("Warn")
	if err != nil {
		return err
	}
	Info, err = initLogger("Info")
	if err != nil {
		return err
	}
	Debug, err = initLogger("Debug")
	if err != nil {
		return err
	}
	return nil
}

func initLogger(name string) (*log.Logger, error) {
	varSym, errSym := plug.Lookup(name)
	if errSym != nil {
		return nil, errSym
	}
	logger, okResolve := varSym.(**log.Logger)
	if !okResolve {
		return nil, errResolve
	}
	return *logger, nil
}
