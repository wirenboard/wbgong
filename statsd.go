package wbgo

import (
	"log"
	"time"

	"github.com/alexcesaro/statsd"
)

var (
	funcNewStatsdClientWrapper    func(string, ...statsd.Option) (StatsdClientWrapper, error)
	funcNewStatsdRuntimeCollector func(StatsdClientWrapper) StatsdRuntimeCollector
)

// StatsdClientWrapper is a wrapper around statsd client
type StatsdClientWrapper interface {
	SetCallback(func(*statsd.Client))
	Clone(string, ...statsd.Option) StatsdClientWrapper
	Start(time.Duration)
	Stop()
}

// StatsdRuntimeCollector is a runtime metrics collector for statsd
type StatsdRuntimeCollector interface {
	Start()
	Stop()
}

// NewStatsdClientWrapper returns new StatsdClientWrapper
func NewStatsdClientWrapper(prefix string, opts ...statsd.Option) (StatsdClientWrapper, error) {
	if funcNewStatsdClientWrapper != nil {
		return funcNewStatsdClientWrapper(prefix, opts...)
	}
	funcSym, errSym := plug.Lookup("NewStatsdClientWrapper")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcNewStatsdClientWrapper, okResolve = funcSym.(func(string, ...statsd.Option) (StatsdClientWrapper, error))
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcNewStatsdClientWrapper(prefix, opts...)
}

// NewStatsdRuntimeCollector returns new StatsdRuntimeCollector
func NewStatsdRuntimeCollector(client StatsdClientWrapper) StatsdRuntimeCollector {
	if funcNewStatsdRuntimeCollector != nil {
		return funcNewStatsdRuntimeCollector(client)
	}
	funcSym, errSym := plug.Lookup("NewStatsdRuntimeCollector")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcNewStatsdRuntimeCollector, okResolve = funcSym.(func(StatsdClientWrapper) StatsdRuntimeCollector)
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcNewStatsdRuntimeCollector(client)
}
