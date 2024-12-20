package wbgong

// ControlError is MQTT control error interface
type ControlError interface {
	Error() string
}

// ControlArgs is a handy way to pass control parameters to constructor
type ControlArgs interface {
	SetDevice(Device) ControlArgs
	SetId(string) ControlArgs
	SetTitle(title Title) ControlArgs
	SetEnumTitles(map[string]Title) ControlArgs
	SetDescription(string) ControlArgs
	SetType(string) ControlArgs
	SetUnits(string) ControlArgs
	SetReadonly(bool) ControlArgs
	SetMax(float64) ControlArgs
	SetMin(float64) ControlArgs
	SetPrecision(float64) ControlArgs
	SetError(ControlError) ControlArgs
	SetOrder(int) ControlArgs
	SetRawValue(string) ControlArgs
	SetValue(interface{}) ControlArgs
	SetDoLoadPrevious(bool) ControlArgs
	// SetLazyInit sets lazyInit flag to control
	// If true - control will not create topic in mqtt before once explicitly set value to this control
	// It also means that storage will be not used to store or restore values at all
	SetLazyInit(bool) ControlArgs

	GetDevice() Device
	GetID() *string
	GetTitle() *Title
	GetEnumTitles() *map[string]Title
	GetDescription() *string
	GetType() *string
	GetUnits() *string
	GetReadonly() *bool
	GetMax() *float64
	GetMin() *float64
	GetPrecision() *float64
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

	// Checks whether control is delete
	IsDeleted() bool

	// Checks whether control has retained value
	// (which is not true for button types or something).
	IsRetained() bool

	// Checks whether control belongs to virtual device
	IsVirtual() bool

	// Checks whether user wants to load previous value for local control
	DoLoadPrevious() bool

	// generic getters
	GetId() string                   // Gets control id (/devices/+/controls/[id])
	GetTitle() Title                 // Gets control title (/meta/title)
	GetEnumTitles() map[string]Title // Gets control enum titles (/meta/enum)
	GetDescription() string          // Gets control description (/meta/description)
	GetType() string                 // Gets control type string (/meta/type) (TODO: special type for this)
	GetUnits() string                // Gets control value units (/meta/units)
	GetReadonly() bool               // Checks whether control is read only
	GetMax() float64                 // Gets max value for 'range'/'value' type (FIXME: rework this?)
	GetMin() float64                 // Gets min value for 'range'/'value' type
	GetPrecision() float64           // Gets precision for 'value' type
	GetError() ControlError          // Gets control error (/meta/error)
	GetOrder() int                   // Gets control order (or -1 for auto) (/meta/order)
	GetValue() (interface{}, error)  // Gets control value (converted according to type)
	GetRawValue() string             // Gets control value string
	GetLazyInit() bool               // Gets control lazyInit flag

	// generic setters
	SetDescription(desc string) FuncError
	SetTitle(title Title) FuncError
	SetEnumTitles(map[string]Title) FuncError
	SetType(t string) FuncError
	SetUnits(units string) FuncError
	SetReadonly(r bool) FuncError
	SetMax(max float64) FuncError
	SetMin(min float64) FuncError
	SetPrecision(prec float64) FuncError
	SetError(e ControlError) FuncError
	SetOrder(ord int) FuncError
	// SetLazyInit sets lazyInit flag to control
	// If true - control will not create topic in mqtt before once explicitly set value to this control
	// It also means that storage will be not used to store or restore values at all
	SetLazyInit(bool) FuncError

	// Updates control value for local device
	// and notifies subscribers if flag is set
	UpdateValue(val interface{}, notifySubs bool) FuncError

	// Sets '/on' value for external devices
	SetOnValue(val interface{}) FuncError

	// Gets all metadata from control (for driver)
	GetMeta() MetaInfo

	// Gets all metadata from control for /meta (for driver)
	GetMetaJson() MetaInfo

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
	AcceptMeta(metaType, value string) error
}

// ControlValueHandler is a function that handles new values on /devices/+/controls/+
// XXX: TBD: reaction on error?
// Handlers are running sync with DriverFrontend, so try not to push heavy stuff here
type ControlValueHandler func(control Control, value interface{}, prevValue interface{}, tx DriverTx) error
