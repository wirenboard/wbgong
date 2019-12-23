package wbgong

// ControlError is MQTT control error interface
type ControlError interface {
	Error() string
}

// ControlArgs is a handy way to pass control parameters to constructor
type ControlArgs interface {
	SetDevice(Device) ControlArgs
	SetId(string) ControlArgs
	SetDescription(string) ControlArgs
	SetType(string) ControlArgs
	SetUnits(string) ControlArgs
	SetReadonly(bool) ControlArgs
	SetWritable(bool) ControlArgs
	SetMax(int) ControlArgs
	SetError(ControlError) ControlArgs
	SetOrder(int) ControlArgs
	SetRawValue(string) ControlArgs
	SetValue(interface{}) ControlArgs
	SetDoLoadPrevious(bool) ControlArgs
	SetLazyInit(bool) ControlArgs

	GetDevice() Device
	GetID() *string
	GetDescription() *string
	GetType() *string
	GetUnits() *string
	GetWritable() *bool
	GetMax() *int
	GetError() ControlError
	GetOrder() *int
	GetRawValue() *string
	GetValue() interface{}
	GetDoLoadPrevious() *bool
	GetLazyInit() *bool
}

// Control is a user representation of MQTT device control
type Control interface {
	// Gets device-owner
	GetDevice() Device

	// Sets transaction object
	SetTx(tx DriverTx)

	// Checks whether control is complete (all required metadata received)
	IsComplete() bool

	// Checks whether control has retained value
	// (which is not true for button types or something).
	IsRetained() bool

	// Checks whether control belongs to virtual device
	IsVirtual() bool

	// Checks whether user wants to load previous value for local control
	DoLoadPrevious() bool

	// generic getters
	GetId() string                  // Gets control id (/devices/+/controls/[id])
	GetDescription() string         // Gets control description (/meta/description)
	GetType() string                // Gets control type string (/meta/type) (TODO: special type for this)
	GetUnits() string               // Gets control value units (/meta/units)
	GetReadonly() bool              // Checks whether control is read only (TODO: merge with Writable?)
	GetWritable() bool              // Checks whether control is writable
	GetMax() int                    // Gets max value for 'range' type (FIXME: rework this?)
	GetError() ControlError         // Gets control error (/meta/error)
	GetOrder() int                  // Gets control order (or -1 for auto) (/meta/order)
	GetValue() (interface{}, error) // Gets control value (converted according to type)
	GetRawValue() string            // Gets control value string
	GetLazyInit() bool              // Gets control lazyInit flag

	// generic setters
	SetDescription(desc string) FuncError
	SetType(t string) FuncError
	SetUnits(units string) FuncError
	SetReadonly(r bool) FuncError
	SetWritable(w bool) FuncError
	SetMax(max int) FuncError
	SetError(e ControlError) FuncError
	SetOrder(ord int) FuncError
	SetLazyInit(bool) FuncError

	// universal interface for UpdateValue and SetOnValue
	SetValue(val interface{}) FuncError

	// Updates control value for local device
	UpdateValue(val interface{}) FuncError

	// Sets '/on' value for external devices
	SetOnValue(val interface{}) FuncError

	// Gets all metadata from control (for driver)
	GetMeta() MetaInfo

	// Saves single meta value in control structure (for driver)
	SetSingleMeta(meta string, value string) error

	// Sets new value handler (for external controls only)
	SetValueUpdateHandler(f ControlValueHandler) error

	// Sets new 'on' value handler (for local controls only)
	SetOnValueReceiveHandler(f ControlValueHandler) error

	// Marks control as deleted
	// Used by LocalDevice in RemoveControl
	MarkDeleted()

	// Sets raw control value
	// This method is used by driver to update values
	// for controls - without notifying
	SetRawValue(value string) error

	//
	// Methodes to accept values from MQTT, called generally by driver
	AcceptValue(rawValue string) error
	AcceptOnValue(rawValue string) error
	AcceptMeta(event NewExternalDeviceControlMetaEvent) error
}

// ControlValueHandler is a function that handles new values on /devices/+/controls/+
// XXX: TBD: reaction on error?
// Handlers are running sync with DriverFrontend, so try not to push heavy stuff here
type ControlValueHandler func(control Control, value interface{}, tx DriverTx) error
