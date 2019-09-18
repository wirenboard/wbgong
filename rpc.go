package wbgo

// MQTTRPCServer represents mqtt rpc server
type MQTTRPCServer interface {
	Start()
	Stop()
	Register(interface{}) error
	RegisterName(string, interface{}) error
}
