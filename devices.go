package wbgong

import (
	"fmt"
	"sort"
)

// Device is a user representation of MQTT device
//
// Devices are passed to observers, but observers must not
// access Devices async with Driver (use them only as references
// to pass them via Events)
//
// External device objects are only created by Driver; driver doesn't access
// to these objects after passing AddExternalDevice event
type Device interface {

	// Gets external driver ID (from /meta/driver)
	GetDriverId() string

	// Gets local driver object
	GetDriver() DeviceDriver

	// Sets transaction
	SetTx(tx DriverTx)

	// Gets device id
	GetId() string

	// Gets device title
	GetTitle() string

	// Gets control by id
	GetControl(id string) Control

	// Lists controls
	ControlsList() []Control

	// Checks whether control with given id exists
	HasControl(id string) bool

	// Gets all device metadata
	GetMeta() MetaInfo

	// Marks device as deleted
	// Used by Driver frontend and device itself
	MarkDeleted()

	// Checks if device is marked as deleted
	IsDeleted() bool
}

// LocalDevice is a user representation of local MQTT device
type LocalDevice interface {
	Device

	// Checks whether device is virtual
	// * Virtual devices have value cache (ability to restore previous control values on load);
	// * Non-virtual devices (physical devices) doesn't
	// This flag is used by Driver
	IsVirtual() bool

	// Checks whether user wants to load previous values for device at load.
	// This flag is used by Driver
	DoLoadPrevious() bool

	// Creates new control by args
	// Returns 'future' value
	CreateControl(args ControlArgs) func() (Control, error)

	// Removes control by id
	RemoveControl(id string) func() error
}

// ExternalDevice is a user representation of external MQTT device
// Provides extra methods to modify device via driver events
type ExternalDevice interface {
	Device

	// Add external control
	AddControl(control Control) error

	// Adds meta information from event
	AcceptMeta(event NewExternalDeviceMetaEvent) error

	// Sets device title
	SetTitle(title string)

	// Internal function to remove control on cleanup
	RemoveControl(id string)
}

// LocalDeviceArgs is a handy way to pass local device attributes
type LocalDeviceArgs interface {
	SetVirtual(v bool) LocalDeviceArgs
	SetDoLoadPrevious(v bool) LocalDeviceArgs
	SetId(v string) LocalDeviceArgs
	SetTitle(v string) LocalDeviceArgs
	SetDriver(v DeviceDriver) LocalDeviceArgs
	GetID() *string
	GetTitle() *string
	GetDriver() DeviceDriver
	GetVirtual() bool
	GetDoLoadPrevious() bool
}

// MetaInfo is a type that represents /meta/+ topics for drivers and controls
type MetaInfo map[string]string

// Implementation of Stringer interface to print metadata correctly
func (m MetaInfo) String() (ret string) {
	ret = "[ "
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		ret += fmt.Sprintf("%s: '%s' ", key, m[key])
	}

	ret += "]"
	return
}

// Delta calculates the difference between two MetaInfos and returns
// a new MetaInfo with delta containing these rows:
//  - rows with the same key but different values;
//  - rows from newMeta without the complementary ones in oldMeta - with values from newMeta;
//  - rows from oldMeta without the complementary ones in newMeta - with empty values ("").
//
//  For example:
//
//  newMeta: [ "a": "123", "b": "xyz", "d": "0" ]
//  oldMeta: [ "a": "456", "c": "hello", "d": "0" ]
//  delta: [ "a": "123", "b": "xyz", "c": "" ]
func (m MetaInfo) Delta(oldMeta MetaInfo) (delta MetaInfo) {
	delta = make(MetaInfo)

	// find changed values
	for key, newValue := range m {
		if oldValue, ok := oldMeta[key]; !ok || (ok && oldValue != newValue) {
			delta[key] = newValue
		}
	}

	// find deleted values
	for key := range oldMeta {
		if _, ok := m[key]; !ok {
			delta[key] = ""
		}
	}

	return
}
