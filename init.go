package wbgong

import (
	"errors"
	"log"
	"plugin"
)

var errResolve = errors.New("Could not resolve symbol")

var plug *plugin.Plugin

// Init tries to load shared library
func Init(path string) (err error) {
	plug, err = plugin.Open(path)
	if err != nil {
		log.Printf("Error: '%s'", err)
		return err
	}
	return nil
}
