// Device filters allows to subscribe only on limited list
// of devices/controls.
// Device filters are used in DriverBackend to filter MQTT topics to subscribe

package wbgong

// DeviceFilter is a device filter representation for Driver
type DeviceFilter interface {
	// Topics() returns list of deviceId/controlId pairs to subscribe to
	Topics() []DeviceControlPair

	// MatchTopic checks whether MQTT topic matches this filter
	MatchTopic(t string) bool
}

// DeviceControlPair is a structure contains deviceId and controlId
// perfect for filtering
type DeviceControlPair struct {
	deviceID  string
	controlID string
}

func (dcp *DeviceControlPair) GetDeviceID() string {
	return dcp.deviceID
}

func (dcp *DeviceControlPair) GetControlID() string {
	return dcp.controlID
}

// AllDevicesFilter is a filter which allows all devices
type AllDevicesFilter struct{}

func (f *AllDevicesFilter) Topics() []DeviceControlPair {
	return []DeviceControlPair{
		DeviceControlPair{CONV_SUBTOPIC_ALL, CONV_SUBTOPIC_ALL},
	}
}

func (f *AllDevicesFilter) MatchTopic(t string) bool {
	return true
}

// NoDevicesFilter is a filter which denies all devices
type NoDevicesFilter struct{}

func (f *NoDevicesFilter) Topics() []DeviceControlPair {
	return []DeviceControlPair{}
}

func (f *NoDevicesFilter) MatchTopic(t string) bool {
	return false
}

// DeviceListFilter is a filter which allows only given list of external devices
type DeviceListFilter struct {
	devices map[string]bool
}

func (f *DeviceListFilter) Topics() []DeviceControlPair {
	r := make([]DeviceControlPair, 0, len(f.devices))

	for dev := range f.devices {
		r = append(r, DeviceControlPair{dev, CONV_SUBTOPIC_ALL})
	}

	return r
}

func (f *DeviceListFilter) MatchTopic(t string) bool {
	// we presume that list of topics was used as is, so all topics from list are allowed
	return true
}

func NewDeviceListFilter(devices ...string) *DeviceListFilter {
	d := &DeviceListFilter{
		devices: make(map[string]bool),
	}

	for _, dev := range devices {
		d.devices[dev] = true
	}

	return d
}
