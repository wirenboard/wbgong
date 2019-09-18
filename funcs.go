package wbgo

import (
	"log"
)

var (
	funcUseSyslog           func()
	funcSetDebuggingEnabled func(bool)
	funcDebuggingEnabled    func() bool
	funcSetDebugLogger      func(logger *log.Logger, keep bool)
	funcEnableMQTTDebugLog  func(bool)
	funcMaybeInitProfiling  func(<-chan struct{})
	funcNewDriverBase       func(args DriverArgs) (DeviceDriver, error)
	funcNewLocalDeviceArgs  func() LocalDeviceArgs
	funcNewControlArgs      func() ControlArgs
	funcNewContentTracker   func() ContentTracker
	funcNewDirWatcher       func(string, DirWatcherClient) DirWatcher
	funcNewMQTTRPCServer    func(string, MQTTClient) MQTTRPCServer
)

// UseSyslog sets logs output to syslog
func UseSyslog() {
	if funcUseSyslog != nil {
		funcUseSyslog()
		return
	}
	funcSym, errSym := plug.Lookup("UseSyslog")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcUseSyslog, okResolve = funcSym.(func())
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	funcUseSyslog()
	return
}

// SetDebuggingEnabled on/off debugging
func SetDebuggingEnabled(enable bool) {
	if funcSetDebuggingEnabled != nil {
		funcSetDebuggingEnabled(enable)
		return
	}
	funcSym, errSym := plug.Lookup("SetDebuggingEnabled")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcSetDebuggingEnabled, okResolve = funcSym.(func(bool))
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	funcSetDebuggingEnabled(enable)
	return
}

// DebuggingEnabled returns true if debugging is enabled
func DebuggingEnabled() bool {
	if funcDebuggingEnabled != nil {
		return funcDebuggingEnabled()
	}
	funcSym, errSym := plug.Lookup("DebuggingEnabled")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcDebuggingEnabled, okResolve = funcSym.(func() bool)
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcDebuggingEnabled()
}

// SetDebugLogger sets prepared debug logger
func SetDebugLogger(logger *log.Logger, keep bool) {
	if funcSetDebugLogger != nil {
		funcSetDebugLogger(logger, keep)
		return
	}
	funcSym, errSym := plug.Lookup("SetDebugLogger")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcSetDebugLogger, okResolve = funcSym.(func(logger *log.Logger, keep bool))
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	funcSetDebugLogger(logger, keep)
	return
}

// EnableMQTTDebugLog enables mqtt debug logging
func EnableMQTTDebugLog(useSyslog bool) {
	if funcEnableMQTTDebugLog != nil {
		funcEnableMQTTDebugLog(useSyslog)
		return
	}
	funcSym, errSym := plug.Lookup("EnableMQTTDebugLog")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcEnableMQTTDebugLog, okResolve = funcSym.(func(bool))
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	funcEnableMQTTDebugLog(useSyslog)
	return
}

// MaybeInitProfiling enables cpu profiling if needed
func MaybeInitProfiling(readyCh <-chan struct{}) {
	if funcMaybeInitProfiling != nil {
		funcMaybeInitProfiling(readyCh)
		return
	}
	funcSym, errSym := plug.Lookup("MaybeInitProfiling")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcMaybeInitProfiling, okResolve = funcSym.(func(<-chan struct{}))
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	funcMaybeInitProfiling(readyCh)
	return
}

// NewDriverBase returns new base driver
func NewDriverBase(args DriverArgs) (DeviceDriver, error) {
	if funcNewDriverBase != nil {
		return funcNewDriverBase(args)
	}
	funcSym, errSym := plug.Lookup("NewDriverBase")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcNewDriverBase, okResolve = funcSym.(func(args DriverArgs) (DeviceDriver, error))
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcNewDriverBase(args)
}

// NewLocalDeviceArgs return new LocalDeviceArgs
func NewLocalDeviceArgs() LocalDeviceArgs {
	if funcNewLocalDeviceArgs != nil {
		return funcNewLocalDeviceArgs()
	}
	funcSym, errSym := plug.Lookup("NewLocalDeviceArgs")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcNewLocalDeviceArgs, okResolve = funcSym.(func() LocalDeviceArgs)
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcNewLocalDeviceArgs()
}

// NewControlArgs return new ControlArgs
func NewControlArgs() ControlArgs {
	if funcNewControlArgs != nil {
		return funcNewControlArgs()
	}
	funcSym, errSym := plug.Lookup("NewControlArgs")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcNewControlArgs, okResolve = funcSym.(func() ControlArgs)
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcNewControlArgs()
}

// NewContentTracker return new ContentTracker
func NewContentTracker() ContentTracker {
	if funcNewContentTracker != nil {
		return funcNewContentTracker()
	}
	funcSym, errSym := plug.Lookup("NewContentTracker")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcNewContentTracker, okResolve = funcSym.(func() ContentTracker)
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcNewContentTracker()
}

// NewDirWatcher return new DirWatcher
func NewDirWatcher(pattern string, client DirWatcherClient) DirWatcher {
	if funcNewDirWatcher != nil {
		return funcNewDirWatcher(pattern, client)
	}
	funcSym, errSym := plug.Lookup("NewDirWatcher")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcNewDirWatcher, okResolve = funcSym.(func(string, DirWatcherClient) DirWatcher)
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcNewDirWatcher(pattern, client)
}

// NewMQTTRPCServer return new MQTTRPCServer
func NewMQTTRPCServer(appName string, mqttClient MQTTClient) MQTTRPCServer {
	if funcNewMQTTRPCServer != nil {
		return funcNewMQTTRPCServer(appName, mqttClient)
	}
	funcSym, errSym := plug.Lookup("NewMQTTRPCServer")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcNewMQTTRPCServer, okResolve = funcSym.(func(string, MQTTClient) MQTTRPCServer)
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcNewMQTTRPCServer(appName, mqttClient)
}
