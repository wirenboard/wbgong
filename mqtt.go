package wbgong

import "log"

var (
	funcNewPahoMQTTClient func(string, string) MQTTClient
)

// MQTTMessageHandler is a handler of MQTTMessages
type MQTTMessageHandler func(message MQTTMessage)

// MQTTMessage represents mqtt message
type MQTTMessage struct {
	Topic    string
	Payload  string
	QoS      byte
	Retained bool
}

// MQTTClient is a mqtt client interface
type MQTTClient interface {
	WaitForRetained(callback func())
	Start()
	Stop()
	Publish(message MQTTMessage)
	Subscribe(callback MQTTMessageHandler, topics ...string)
	Unsubscribe(topics ...string)
}

// NewPahoMQTTClient returns new Paho mqtt client
func NewPahoMQTTClient(server, clientID string) MQTTClient {
	if funcNewPahoMQTTClient != nil {
		return funcNewPahoMQTTClient(server, clientID)
	}
	funcSym, errSym := plug.Lookup("NewPahoMQTTClient")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %s", errSym)
	}
	var okResolve bool
	funcNewPahoMQTTClient, okResolve = funcSym.(func(string, string) MQTTClient)
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return funcNewPahoMQTTClient(server, clientID)
}
