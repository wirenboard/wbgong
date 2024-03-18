package testutils

import (
	"testing"
	"time"

	"github.com/wirenboard/wbgong"
)

// FakeDriverBackend is a dummy backend for model testing
// Implements DriverBackend
type FakeDriverBackend struct {
	*Recorder
	externalDeviceFactory wbgong.ExternalDeviceFactory
	controlFactory        wbgong.ControlFactory
	frontend              wbgong.DriverFrontend
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

func (backend *FakeDriverBackend) SetExternalDeviceFactory(e wbgong.ExternalDeviceFactory) {
	backend.externalDeviceFactory = e
}

func (backend *FakeDriverBackend) SetControlFactory(f wbgong.ControlFactory) {
	backend.controlFactory = f
}

func (backend *FakeDriverBackend) SetFrontend(f wbgong.DriverFrontend) {
	backend.frontend = f
}

func (backend *FakeDriverBackend) SetExternalDeviceFilter(e wbgong.DeviceFilter) {}

func closedChan() <-chan error {
	e := make(chan error, 1)
	e <- nil
	return e
}

func (backend *FakeDriverBackend) SetFilter(f wbgong.DeviceFilter) <-chan struct{} {
	backend.Rec("[FakeDriverBackend] SetFilter(%T)", f)
	c := make(chan struct{}, 1)
	c <- struct{}{}
	return c
}

func (backend *FakeDriverBackend) RemoveDevice(dev wbgong.LocalDevice) <-chan error {
	backend.Rec("[FakeDriverBackend] RemoveDevice(%s)", dev.GetId())
	return closedChan()
}

func (backend *FakeDriverBackend) NewDevice(dev wbgong.LocalDevice) <-chan error {
	backend.Rec("[FakeDriverBackend] NewDevice(%s, meta %v)", dev.GetId(), dev.GetMeta())
	return closedChan()
}

func (backend *FakeDriverBackend) NewDeviceControl(control wbgong.Control) <-chan error {
	backend.Rec("[FakeDriverBackend] NewDeviceControl(control %s/%s, value %s, meta %v)", control.GetDevice().GetId(), control.GetId(), control.GetRawValue(), control.GetMeta())
	return closedChan()
}

func (backend *FakeDriverBackend) UpdateControlValue(control wbgong.Control, rawValue string) <-chan error {
	backend.Rec("[FakeDriverBackend] UpdateControlValue(control %s/%s, value %s)", control.GetDevice().GetId(), control.GetId(), rawValue)
	return closedChan()
}

func (backend *FakeDriverBackend) SetOnValue(control wbgong.Control, rawValue string) <-chan error {
	backend.Rec("[FakeDriverBackend] SetOnValue(control %s/%s, value %s)", control.GetDevice().GetId(), control.GetId(), rawValue)
	return closedChan()
}

func (backend *FakeDriverBackend) UpdateDeviceMeta(dev wbgong.LocalDevice, meta string, value interface{}) <-chan error {
	backend.Rec("[FakeDriverBackend] UpdateDeviceMeta(device %s, meta %s, value %s)", dev.GetId(), meta, value)
	return closedChan()
}

func (backend *FakeDriverBackend) UpdateDeviceMetaJson(dev wbgong.LocalDevice) <-chan error {
	backend.Rec("[FakeDriverBackend] UpdateDeviceMetaJson(device %s)", dev.GetId())
	return closedChan()
}

func (backend *FakeDriverBackend) UpdateControlMeta(control wbgong.Control, meta string, value interface{}) <-chan error {
	backend.Rec("[FakeDriverBackend] UpdateControlMeta(control %s/%s, meta %s, value %s)", control.GetDevice().GetId(), control.GetId(), meta, value)
	return closedChan()
}

func (backend *FakeDriverBackend) UpdateControlMetaJson(control wbgong.Control) <-chan error {
	backend.Rec("[FakeDriverBackend] UpdateControlMetaJson(control %s/%s)", control.GetDevice().GetId(), control.GetId())
	return closedChan()
}

func (backend *FakeDriverBackend) RemoveControl(ctrl wbgong.Control) <-chan error {
	backend.Rec("[FakeDriverBackend] RemoveControl(control %s/%s)", ctrl.GetDevice().GetId(), ctrl.GetId())
	return closedChan()
}

func (backend *FakeDriverBackend) RemoveExternalDevice(dev wbgong.ExternalDevice) {
	backend.Rec("[FakeDriverBackend] RemoveExternalDevice(%s)", dev.GetId())
}

func (backend *FakeDriverBackend) RemoveExternalControl(ctrl wbgong.Control) {
	backend.Rec("[FakeDriverBackend] RemoveExternalControl(control %s/%s)", ctrl.GetDevice().GetId(), ctrl.GetId())
}

func (backend *FakeDriverBackend) PushEvent(event wbgong.DriverEvent) {
	backend.frontend.PushEvent(event)
}

// FakeDriverFrontend is a dummy frontend for model testing
// Implements Driver, DriverFrontend and DeviceDriver
type FakeDriverFrontend struct {
	*Recorder

	id      string
	backend wbgong.DriverBackend
	doReown bool

	Devices map[string]wbgong.Device
}

func NewFakeDriverFrontend(id string, t *testing.T) *FakeDriverFrontend {
	return &FakeDriverFrontend{
		Recorder: NewRecorder(t),
		id:       id,
		Devices:  make(map[string]wbgong.Device),
	}
}

func (f *FakeDriverFrontend) GetId() string {
	return f.id
}

func (f *FakeDriverFrontend) SetFilter(fl wbgong.DeviceFilter) {
	<-f.backend.SetFilter(fl)
}

// dummy
func (f *FakeDriverFrontend) SetLocalDeviceFactory(ff wbgong.LocalDeviceFactory) wbgong.LocalDeviceFactory {
	return nil
}

// dummy
func (f *FakeDriverFrontend) SetExternalDeviceFactory(ff wbgong.ExternalDeviceFactory) wbgong.ExternalDeviceFactory {
	return nil
}

// dummy
func (f *FakeDriverFrontend) SetControlFactory(ff wbgong.ControlFactory) wbgong.ControlFactory {
	return nil
}

// dummy
func (f *FakeDriverFrontend) OnRetainReady(ff func(tx wbgong.DriverTx)) {}

// dummy
func (f *FakeDriverFrontend) WaitForReady() {}

// dummy
func (f *FakeDriverFrontend) OnDriverEvent(ff func(e wbgong.DriverEvent)) wbgong.HandlerID {
	return 0
}

// dummy
func (f *FakeDriverFrontend) RemoveOnDriverEventHandler(i wbgong.HandlerID) {}

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
func (f *FakeDriverFrontend) Access(thunk func(tx wbgong.DriverTx) error) error {
	return nil
}

// dummy
func (f *FakeDriverFrontend) AccessAsync(thunk func(tx wbgong.DriverTx) error) func() error {
	return func() error {
		return nil
	}
}

// dummy
func (f *FakeDriverFrontend) BeginTx() (wbgong.DriverTx, error) {
	return nil, nil
}

// dummy
func (f *FakeDriverFrontend) CreateUnsafeTx() wbgong.DeviceDriverTx {
	return nil
}

// dummy
func (f *FakeDriverFrontend) Close() {}

// Push driver event into processing queue
func (f *FakeDriverFrontend) PushEvent(event wbgong.DriverEvent) {
	f.Rec("[FakeDriverFrontend] Received event %v", event)

	switch e := event.(type) {
	case wbgong.NewExternalDeviceEvent:
		f.Devices[e.Device.GetId()] = e.Device
	case wbgong.NewExternalDeviceControlEvent:
		e.Device.AddControl(e.Control)
	}
}

// Connect backend
func (f *FakeDriverFrontend) SetBackend(b wbgong.DriverBackend) {
	f.backend = b
}

func (f *FakeDriverFrontend) NeedToReownUnknownDevices() bool {
	return f.doReown
}
