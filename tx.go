package wbgong

// DriverTx is a transaction interface
type DriverTx interface {
	// End() closes transaction and makes it invalid
	End()

	// GetDevicesList returns a slice of devices
	// currently registered in driver
	GetDevicesList() []Device

	// CreateDevice creates local device with given attributes
	// CreateDevice is a "future" request, so it returns function
	// which blocks until all required data is received
	CreateDevice(args LocalDeviceArgs) func() (LocalDevice, error)

	// GetDevice gets device by ID
	GetDevice(id string) Device

	// HasDevice checks whether device with given ID exists
	HasDevice(id string) bool

	// RemoveDevice removes local device by reference
	// Device object will be marked as 'deleted' and so all operations
	// with device via references will fail.
	//
	// Returns nil if device is local and removed successfully
	RemoveDevice(dev LocalDevice) func() error

	// RemoveDeviceById tries to find device by its ID
	// and removes it
	//
	// See also RemoveDevice
	RemoveDeviceById(id string) func() error

	// Returns DeviceDriverTx from this Tx
	ToDeviceDriverTx() DeviceDriverTx
}

// DeviceDriverTx extends DriverTx with
// additional methodes to be used by Devices and Controls
type DeviceDriverTx interface {
	DriverTx

	// CreateControl creates control with given attributes
	// Device which creates control must add itself into arguments list
	// CreateControl is a "future" request
	CreateControl(args ControlArgs) func() (Control, error)

	// UpdateControlValue updates value for given control
	// (sends message to /d/+/c/+) and notifies local subscribers
	// about this change if notification flag is set
	UpdateControlValue(control Control, rawValue string, prevRawValue string, notifySubs bool) func() error

	UpdateControlMeta(control Control, meta string, value interface{}) func() error
	UpdateControlMetaJson(control Control) func() error
	UpdateDeviceMeta(dev LocalDevice, meta string, value interface{}) func() error

	// SetOnValue sends /on value for given control
	// SetOnValue is a future request
	SetOnValue(control Control, rawValue string) func() error

	// RemoveControl removes control from device
	// and performs MQTT topics cleanup
	RemoveControl(control Control) func() error

	// RemoveExternalControl removes control from external device
	// on control's request (when all MQTT topics for it are cleared)
	RemoveExternalControl(control Control)

	// RemoveExternalDevice removes external device on device's request
	// (when all MQTT topics are cleared)
	RemoveExternalDevice(dev ExternalDevice)
}
