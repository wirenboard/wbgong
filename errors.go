package wbgong

import "errors"

// FuncError is a basic error-returning promise function
// This function is returned by API methodes for user to control
// operation status asyncronously
type FuncError func() error

// MakeFuncError is a basic FuncError implementation which returns an error immediately
func MakeFuncError(err error) FuncError {
	return func() error {
		return err
	}
}

// EmptyFuncError is a empty func error which return no error
var EmptyFuncError FuncError = func() error {
	return nil
}

// Common errors definitions
var (
	EventQueueFullError       = errors.New("Driver event queue is full")
	DriverActiveError         = errors.New("Driver loop is already running")
	DriverInactiveError       = errors.New("Driver loop is not running")
	DriverWrongArgumentsError = errors.New("Wrong arguments set for NewDriver")
	DriverTimeoutError        = errors.New("Driver timeout error")
	DeviceRedefinitionError   = errors.New("Device redefinition")
	ControlRedefinitionError  = errors.New("Control redefinition")
	NonLocalControlError      = errors.New("Trying to register non-local control")
	DeviceIdMissingError      = errors.New("Device ID is missing")
	DeviceAlreadyExistsError  = errors.New("Device with given ID already exists")
	DeviceNotExistError       = errors.New("Device with given ID doesn't exist")
	ControlArgsMissingError   = errors.New("Some arguments for CreateControl missing")
	IncorrectDeviceIdError    = errors.New("Device ID has incorrect symbols")
	StorageUnavailableError   = errors.New("External storage is not initialized")
	StorageValueNotFoundError = errors.New("No value in storage")

	LocalDeviceError    = errors.New("Device is local")
	ExternalDeviceError = errors.New("Device is external")
	BackendActiveError  = errors.New("Driver backend is running already")

	UnknownDeviceMetaError    = errors.New("Unknown device meta type")
	ControlAlreadyExistsError = errors.New("Control already exists")
	NoSuchControlError        = errors.New("No such control")
	LocalDeviceArgumentsError = errors.New("Wrong local device arguments list, check required fields")
	ControlArgumentsError     = errors.New("Wrong control arguments list, check required fields")
	DeviceDeletedError        = errors.New("This device was deleted")

	ExternalControlError    = errors.New("This control is external")
	LocalControlError       = errors.New("This control is local")
	WrongValueTypeError     = errors.New("Wrong value type")
	UnknownControlMetaError = errors.New("Unknown control meta type")
	IncompleteControlError  = errors.New("This control is incomplete")
	ControlDeletedError     = errors.New("This control was deleted")
	IncorrectControlIdError = errors.New("Control ID is incorrect")
	NoTxContextError        = errors.New("No Tx context")
	NotWritableControlError = errors.New("This control is not writable")
	ReadonlyMissingError    = errors.New("Missing of mandatory readonly argument")
)
