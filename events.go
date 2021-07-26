package wbgong

import "fmt"

// StartEvent represents start event
type StartEvent struct{}

func (e StartEvent) String() string {
	return "StartEvent{}"
}

// StopEvent represents stop event
type StopEvent struct{}

func (e StopEvent) String() string {
	return "StopEvent{}"
}

// ReadyEvent is a driver ready event
// Means that all retained messages are received
type ReadyEvent struct{}

func (e ReadyEvent) String() string {
	return "ReadyEvent{}"
}

// NewExternalDeviceEvent event of external device is detected
// Fires when driver receives first message with matching topic
// Passes new ExternalDevice object representing this device
type NewExternalDeviceEvent struct {
	Device ExternalDevice
}

func (e NewExternalDeviceEvent) String() string {
	return fmt.Sprintf("NewExternalDeviceEvent{Device:%s}", e.Device.GetId())
}

// NewExternalDeviceControlEvent a new external device control is detected
// Fires when driver receives first message with matching topic
// Passes new Control object
type NewExternalDeviceControlEvent struct {
	Device  ExternalDevice
	Control Control
}

func (e NewExternalDeviceControlEvent) String() string {
	return fmt.Sprintf("NewExternalDeviceControlEvent{Device:%s,Control:%s}", e.Device.GetId(), e.Control.GetId())
}

// NewExternalDeviceMetaEvent a new external device meta received
// Fires when driver receives message with topic matching to some
// external device's meta
type NewExternalDeviceMetaEvent struct {
	Device ExternalDevice
	Type   string
	Value  string
}

func (e NewExternalDeviceMetaEvent) String() string {
	return fmt.Sprintf("NewExternalDeviceMetaEvent{Device:%s,Type:%s,Value:%s}", e.Device.GetId(), e.Type, e.Value)
}

// NewExternalDeviceControlMetaEvent a new external device control metadata received
type NewExternalDeviceControlMetaEvent struct {
	Control Control
	Type    string
	Value   string
}

func (e NewExternalDeviceControlMetaEvent) String() string {
	return fmt.Sprintf("NewExternalDeviceControlMetaEvent{Device:%s,Control:%s,Type:%s,Value:%s}",
		e.Control.GetDevice().GetId(), e.Control.GetId(), e.Type, e.Value)
}

// ControlValueEvent a new device value received
// Device may be either local or external. For local controls
// this event ignored by driver (sent for user handlers only)
type ControlValueEvent struct {
	Control      Control
	RawValue     string
	PrevRawValue string
}

func (e ControlValueEvent) String() string {
	return fmt.Sprintf("ControlValueEvent{Device:%s,Control:%s,Value:%s->%s}", e.Control.GetDevice().GetId(),
		e.Control.GetId(), e.PrevRawValue, e.RawValue)
}

// ControlOnValueEvent control received 'on' value
// Valid for local devices only
type ControlOnValueEvent struct {
	Control  Control
	RawValue string
}

func (e ControlOnValueEvent) String() string {
	return fmt.Sprintf("ControlOnValueEvent{Device:%s,Control:%s,Value:%s}", e.Control.GetDevice().GetId(),
		e.Control.GetId(), e.RawValue)
}
