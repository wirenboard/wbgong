package wbgo

import (
	"log"
	"os"
	"time"
)

var (
	funcNewDriverArgs func() DriverArgs
)

// HandlerID is an index of specific event handler
type HandlerID int

// DriverArgs is a driver args
type DriverArgs interface {
	SetId(id string) DriverArgs
	SetMqtt(mqtt MQTTClient) DriverArgs
	SetTesting() DriverArgs
	SetBackend(b DriverBackend) DriverArgs
	SetUseStorage(v bool) DriverArgs
	SetStoragePath(path string) DriverArgs
	SetStorageMode(mode os.FileMode) DriverArgs
	SetReownUnknownDevices(v bool) DriverArgs
	SetStatsdClient(c StatsdClientWrapper) DriverArgs
	Finalize()
	GetBackend() DriverBackend
	GetID() string
	GetIsTesting() bool
	GetReownUnknownDevices() bool
	GetStatsdClient() StatsdClientWrapper
	GetUseStorage() bool
	GetStoragePath() string
	GetStorageMode() os.FileMode
}

// DriverBackend is a backend interface for Driver
// Driver can send requests to Backend via this interface
//
// DriverBackend must be ready to work sync with Driver (at this entrypoints)
type DriverBackend interface {
	Start() error
	Stop()

	SetFrontend(f DriverFrontend)
	SetExternalDeviceFactory(f ExternalDeviceFactory)
	SetControlFactory(f ControlFactory)

	SetFilter(f DeviceFilter) <-chan struct{}
	RemoveDevice(dev LocalDevice) <-chan error
	RemoveControl(ctrl Control) <-chan error
	NewDevice(dev LocalDevice) <-chan error
	NewDeviceControl(control Control) <-chan error
	UpdateControlValue(control Control, rawValue string) <-chan error
	SetOnValue(control Control, rawValue string) <-chan error

	UpdateControlMeta(control Control, meta, value string) <-chan error
	UpdateDeviceMeta(dev LocalDevice, meta, value string) <-chan error

	// these are sent by suicide devices (when all device/control info is cleared)
	RemoveExternalDevice(dev ExternalDevice)
	RemoveExternalControl(ctrl Control)
}

// DriverFrontend is an object DriverBackend interacts with
// DriverFrontend represents Driver as part of FB pair, providing
// special methodes for DriverBackend
type DriverFrontend interface {
	// Gets driver ID
	GetId() string

	// Push driver event into processing queue
	// If block parameter is true, function will block until event is pushed
	// Otherwise, is event queue is full, PushEvent will return EventQueueFullErrror
	PushEvent(event DriverEvent)

	// Connect backend
	SetBackend(b DriverBackend)
}

// Driver is an object user interacts with.
//
// Driver is connected to DriverBackend in given way:
// -> from Driver to Driver - PushEvent()
// -> from Driver to Driver - Backend
//
// Driver stores Devices - user's representation of MQTT devices
// User can interact with them synchronously with driver
type Driver interface {
	//
	// Userspace methods

	//
	// Thread-safe methods
	//

	// GetId returns driver ID from MQTT (published in /devices/+/meta/driver)
	GetId() string

	// SetFilter sets external device filter (NoDevices set by default to omit
	// all external devices).
	// After SetFilter call, backend will reload retained messages, so Ready event will
	// be emitted after
	SetFilter(filter DeviceFilter)

	// LoopOnce tries to receive an event from event queue and process it
	// It quits if quit signal received or timeout occures.
	//
	// It quit signal received, LoopOnce returns false, true otherwise
	//
	// User-defined event handlers also run in loop
	LoopOnce(timeout time.Duration) bool

	// StartLoop
	StartLoop() error

	StopLoop() error

	// BeginTx begins transaction - creates transaction object
	// and locks driver thread to synchronize data access.
	//
	// Don't forget to close transaction after you're done
	// (by using tx.End())
	BeginTx() (DriverTx, error)

	// Access executes user function with transaction and automatically
	// closes it after it's done. Driver loop is blocked in order to provide
	// exclusive access for reader
	//
	// If user function returns error, Access will pass it through.
	//
	// Note that DriverTx object here is only valid within the function.
	Access(thunk func(tx DriverTx) error) error

	// AccessAsync works like Access but executes given function in driver loop
	// and returns Future which waits for given function to return and
	// passes error through.
	AccessAsync(thunk func(tx DriverTx) error) func() error

	// OnRetainReady executes given function in driver loop if
	// all retained messages are received (after ReadyEvent{})
	// To react on each ReadyEvent{}, you need to register function again and again.
	// Handler runs in frontend loop, so all devices access operations are safe.
	OnRetainReady(thunk func(tx DriverTx))

	// WaitForReady locks until next ReadyEvent received
	WaitForReady()

	// OnDriverEvent allows user application to handle backend events.
	// Handler runs in frontend loop, so all devices access operations are safe.
	OnDriverEvent(handler func(e DriverEvent)) HandlerID

	// Removes OnDriverEvent handler
	RemoveOnDriverEventHandler(handlerID HandlerID)

	// Close closes all opened files and connections
	Close()

	//
	// Service methodes

	// SetLocalDeviceFactory registers a function which creates local devices
	SetLocalDeviceFactory(f LocalDeviceFactory) LocalDeviceFactory

	// SetExternalDeviceFactory registers a function which creates external devices
	SetExternalDeviceFactory(f ExternalDeviceFactory) ExternalDeviceFactory

	// SetControlFactory registers a function which creates controls
	SetControlFactory(f ControlFactory) ControlFactory
}

// DeviceDriver extends Driver to be used in transactions
type DeviceDriver interface {
	Driver

	// Creates unsafe transaction to be used in sync handlers
	CreateUnsafeTx() DeviceDriverTx

	// Checks whether we need to reown devices with empty /meta/driver
	NeedToReownUnknownDevices() bool
}

// DriverEvent is a driver event representation for Driver
// Driver sends events to Drivers via Driver.PushEvent(event)
// Implements Stringer
type DriverEvent interface {
	String() string
}

// ExternalDeviceFactory is a function which creates external devices
type ExternalDeviceFactory func(id string, driver DeviceDriver) (ExternalDevice, error)

// ControlFactory is a function which creates controls for devices
// Control checks automatically if it belongs to local or external device
// and modifies its behaviour though
type ControlFactory func(args ControlArgs) (Control, error)

// LocalDeviceFactory is a function which creates local devices with given arguments
type LocalDeviceFactory func(args LocalDeviceArgs) (LocalDevice, error)

// NewDriverArgs returns new driver arguments
func NewDriverArgs() DriverArgs {
	if funcNewDriverArgs != nil {
		return funcNewDriverArgs()
	}
	funcSym, errSym := plug.Lookup("NewDriverArgs")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcNewDriverArgs, okResolve = funcSym.(func() DriverArgs)
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcNewDriverArgs()
}
