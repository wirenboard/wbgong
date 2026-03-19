package wbgong

// MQTTRPCServer represents mqtt rpc server
type MQTTRPCServer interface {
	Start()
	Stop()
	Register(any) error
	RegisterName(string, any) error
}
