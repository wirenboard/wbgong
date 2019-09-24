package testutils

import (
	"testing"
	"time"

	"github.com/contactless/wbgo"
)

// FakeDriverBackend is a dummy backend for model testing
// Implements DriverBackend
type FakeDriverBackend struct {
	*Recorder
	externalDeviceFactory wbgo.ExternalDeviceFactory
	controlFactory        wbgo.ControlFactory
	frontend              wbgo.DriverFrontend
}

func NewFakeDriverBackend(t *testing.T) *FakeDriverBackend {
	backend := &FakeDriverBackend{
		Recorder: NewRecorder(t),
	}

	return backend
}

func (backend *FakeDriverBackend) Start() error {
	backend.Rec("[FakeDriverBackend] Start()")
	return nil
}

func (backend *FakeDriverBackend) Stop() {
	backend.Rec("[FakeDriverBackend] Stop()")
}

func (backend *FakeDriverBackend) SetExternalDeviceFactory(e wbgo.ExternalDeviceFactory) {
	backend.externalDeviceFactory = e
}

func (backend *FakeDriverBackend) SetControlFactory(f wbgo.ControlFactory) {
	backend.controlFactory = f
}

func (backend *FakeDriverBackend) SetFrontend(f wbgo.DriverFrontend) {
	backend.frontend = f
}

func (backend *FakeDriverBackend) SetExternalDeviceFilter(e wbgo.DeviceFilter) {}

func closedChan() <-chan error {
	e := make(chan error, 1)
	e <- nil
	return e
}

func (backend *FakeDriverBackend) SetFilter(f wbgo.DeviceFilter) <-chan struct{} {
	backend.Rec("[FakeDriverBackend] SetFilter(%T)", f)
	c := make(chan struct{}, 1)
	c <- struct{}{}
	return c
}

func (backend *FakeDriverBackend) RemoveDevice(dev wbgo.LocalDevice) <-chan error {
	backend.Rec("[FakeDriverBackend] RemoveDevice(%s)", dev.GetId())
	return closedChan()
}

func (backend *FakeDriverBackend) NewDevice(dev wbgo.LocalDevice) <-chan error {
	backend.Rec("[FakeDriverBackend] NewDevice(%s, meta %v)", dev.GetId(), dev.GetMeta())
	return closedChan()
}

func (backend *FakeDriverBackend) NewDeviceControl(control wbgo.Control) <-chan error {
	backend.Rec("[FakeDriverBackend] NewDeviceControl(control %s/%s, value %s, meta %v)", control.GetDevice().GetId(), control.GetId(), control.GetRawValue(), control.GetMeta())
	return closedChan()
}

func (backend *FakeDriverBackend) UpdateControlValue(control wbgo.Control, rawValue string) <-chan error {
	backend.Rec("[FakeDriverBackend] UpdateControlValue(control %s/%s, value %s)", control.GetDevice().GetId(), control.GetId(), rawValue)
	return closedChan()
}

func (backend *FakeDriverBackend) SetOnValue(control wbgo.Control, rawValue string) <-chan error {
	backend.Rec("[FakeDriverBackend] SetOnValue(control %s/%s, value %s)", control.GetDevice().GetId(), control.GetId(), rawValue)
	return closedChan()
}

func (backend *FakeDriverBackend) UpdateDeviceMeta(dev wbgo.LocalDevice, meta, value string) <-chan error {
	backend.Rec("[FakeDriverBackend] UpdateDeviceMeta(device %s, meta %s, value %s)", dev.GetId(), meta, value)
	return closedChan()
}

func (backend *FakeDriverBackend) UpdateControlMeta(control wbgo.Control, meta, value string) <-chan error {
	backend.Rec("[FakeDriverBackend] UpdateControlMeta(control %s/%s, meta %s, value %s)", control.GetDevice().GetId(), control.GetId(), meta, value)
	return closedChan()
}

func (backend *FakeDriverBackend) RemoveControl(ctrl wbgo.Control) <-chan error {
	backend.Rec("[FakeDriverBackend] RemoveControl(control %s/%s)", ctrl.GetDevice().GetId(), ctrl.GetId())
	return closedChan()
}

func (backend *FakeDriverBackend) RemoveExternalDevice(dev wbgo.ExternalDevice) {
	backend.Rec("[FakeDriverBackend] RemoveExternalDevice(%s)", dev.GetId())
}

func (backend *FakeDriverBackend) RemoveExternalControl(ctrl wbgo.Control) {
	backend.Rec("[FakeDriverBackend] RemoveExternalControl(control %s/%s)", ctrl.GetDevice().GetId(), ctrl.GetId())
}

func (backend *FakeDriverBackend) PushEvent(event wbgo.DriverEvent) {
	backend.frontend.PushEvent(event)
}

// FakeDriverFrontend is a dummy frontend for model testing
// Implements Driver, DriverFrontend and DeviceDriver
type FakeDriverFrontend struct {
	*Recorder

	id      string
	backend wbgo.DriverBackend
	doReown bool

	Devices map[string]wbgo.Device
}

func NewFakeDriverFrontend(id string, t *testing.T) *FakeDriverFrontend {
	return &FakeDriverFrontend{
		Recorder: NewRecorder(t),
		id:       id,
		Devices:  make(map[string]wbgo.Device),
	}
}

func (f *FakeDriverFrontend) GetId() string {
	return f.id
}

func (f *FakeDriverFrontend) SetFilter(fl wbgo.DeviceFilter) {
	<-f.backend.SetFilter(fl)
}

// dummy
func (f *FakeDriverFrontend) SetLocalDeviceFactory(ff wbgo.LocalDeviceFactory) wbgo.LocalDeviceFactory {
	return nil
}

// dummy
func (f *FakeDriverFrontend) SetExternalDeviceFactory(ff wbgo.ExternalDeviceFactory) wbgo.ExternalDeviceFactory {
	return nil
}

// dummy
func (f *FakeDriverFrontend) SetControlFactory(ff wbgo.ControlFactory) wbgo.ControlFactory {
	return nil
}

// dummy
func (f *FakeDriverFrontend) OnRetainReady(ff func(tx wbgo.DriverTx)) {}

// dummy
func (f *FakeDriverFrontend) WaitForReady() {}

// dummy
func (f *FakeDriverFrontend) OnDriverEvent(ff func(e wbgo.DriverEvent)) wbgo.HandlerID {
	return 0
}

// dummy
func (f *FakeDriverFrontend) RemoveOnDriverEventHandler(i wbgo.HandlerID) {}

// dummy
func (f *FakeDriverFrontend) LoopOnce(timeout time.Duration) bool {
	return false
}

// dummy
func (f *FakeDriverFrontend) StartLoop() error {
	return nil
}

func (f *FakeDriverFrontend) StopLoop() error {
	return nil
}

// dummy
func (f *FakeDriverFrontend) Access(thunk func(tx wbgo.DriverTx) error) error {
	return nil
}

// dummy
func (f *FakeDriverFrontend) AccessAsync(thunk func(tx wbgo.DriverTx) error) func() error {
	return func() error {
		return nil
	}
}

// dummy
func (f *FakeDriverFrontend) BeginTx() (wbgo.DriverTx, error) {
	return nil, nil
}

// dummy
func (f *FakeDriverFrontend) CreateUnsafeTx() wbgo.DeviceDriverTx {
	return nil
}

// dummy
func (f *FakeDriverFrontend) Close() {}

// Push driver event into processing queue
func (f *FakeDriverFrontend) PushEvent(event wbgo.DriverEvent) {
	f.Rec("[FakeDriverFrontend] Received event %v", event)

	switch e := event.(type) {
	case wbgo.NewExternalDeviceEvent:
		f.Devices[e.Device.GetId()] = e.Device
	case wbgo.NewExternalDeviceControlEvent:
		e.Device.AddControl(e.Control)
	}
}

// Connect backend
func (f *FakeDriverFrontend) SetBackend(b wbgo.DriverBackend) {
	f.backend = b
}

func (f *FakeDriverFrontend) NeedToReownUnknownDevices() bool {
	return f.doReown
}
