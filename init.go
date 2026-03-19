package wbgong

import (
	"fmt"
	"log"
	"plugin"
)

var plug *plugin.Plugin

// Init tries to load shared library
func Init(path string) (err error) {
	plug, err = plugin.Open(path)
	if err != nil {
		log.Printf("Error: '%s'", err)
		return fmt.Errorf("failed to open plugin: %w", err)
	}
	return nil
}
