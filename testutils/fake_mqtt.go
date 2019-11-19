package testutils

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time" // for Shuffle seed

	"github.com/contactless/wbgong"
)

const (
	DISPATHED_MESSAGE_QUEUE_LEN = 1024 // TODO: make Subscribe easier and set this to 1
)

func topicPartsMatch(pattern []string, topic []string) bool {
	if len(pattern) == 0 {
		return len(topic) == 0
	}

	if pattern[0] == "#" {
		return true
	}

	return len(topic) > 0 &&
		(pattern[0] == "+" || (pattern[0] == topic[0])) &&
		topicPartsMatch(pattern[1:], topic[1:])
}

func topicMatch(pattern string, topic string) bool {
	return topicPartsMatch(strings.Split(pattern, "/"), strings.Split(topic, "/"))
}

func FormatMQTTMessage(message wbgong.MQTTMessage) string {
	suffix := ""
	if message.Retained {
		suffix = ", retained"
	}
	return fmt.Sprintf("[%s] (QoS %d%s)",
		string(message.Payload), message.QoS, suffix)
}

type SubscriptionList []*FakeMQTTClient
type SubscriptionMap map[string]SubscriptionList

type dispatchedMessage struct {
	client  *FakeMQTTClient
	message wbgong.MQTTMessage
	hack    func()
}

type FakeMQTTBroker struct {
	*Recorder
	sync.Mutex
	subscriptions   SubscriptionMap
	waitForRetained bool
	readyChannels   []chan struct{}
	retained        map[string]wbgong.MQTTMessage
	msgQueue        chan dispatchedMessage
	quitCh          chan chan struct{}
	numClients      int
}

func NewFakeMQTTBroker(t *testing.T, rec *Recorder) (broker *FakeMQTTBroker) {
	if rec == nil {
		rec = NewRecorder(t)
	}
	broker = &FakeMQTTBroker{
		Recorder:      rec,
		subscriptions: make(SubscriptionMap),
		retained:      make(map[string]wbgong.MQTTMessage),
		msgQueue:      make(chan dispatchedMessage, DISPATHED_MESSAGE_QUEUE_LEN),
		quitCh:        make(chan chan struct{}, 1),
		numClients:    0,
	}
	return
}

func (broker *FakeMQTTBroker) Start() {
	// start message dispatcher
	go func() {
		for {
			select {
			case msg, ok := <-broker.msgQueue:
				if ok {
					if msg.hack != nil {
						msg.hack()
					} else {
						msg.client.receive(msg.message)
					}
				} else {
					return
				}
			}
		}
	}()
}

func (broker *FakeMQTTBroker) Stop() {
	close(broker.msgQueue)
}

func (broker *FakeMQTTBroker) PushRetainHack(c func()) {
	broker.msgQueue <- dispatchedMessage{hack: c}
}

func (broker *FakeMQTTBroker) SetWaitForRetained(waitForRetained bool) {
	broker.waitForRetained = waitForRetained
}

func (broker *FakeMQTTBroker) SetReady() {
	for _, ch := range broker.readyChannels {
		close(ch)
	}
	broker.readyChannels = nil
}

func (broker *FakeMQTTBroker) Publish(origin string, message wbgong.MQTTMessage) {
	broker.Lock()
	defer broker.Unlock()

	if message.Retained {
		broker.retained[message.Topic] = message
	}

	broker.publish(origin, message)
}

func (broker *FakeMQTTBroker) publish(origin string, message wbgong.MQTTMessage) {
	broker.Rec("%s -> %s: %s", origin, message.Topic, FormatMQTTMessage(message))
	message.Retained = false

	clientsServed := make(map[*FakeMQTTClient]bool)

	for pattern, subs := range broker.subscriptions {
		if !topicMatch(pattern, message.Topic) {
			continue
		}
		for _, client := range subs {
			if clientsServed[client] {
				continue
			}
			broker.queueMessage(client, message)
			clientsServed[client] = true
		}
	}
}

func (broker *FakeMQTTBroker) queueMessage(client *FakeMQTTClient, message wbgong.MQTTMessage) {
	broker.msgQueue <- dispatchedMessage{client, message, nil}
}

func (broker *FakeMQTTBroker) Subscribe(client *FakeMQTTClient, topic string) {
	broker.Lock()
	defer broker.Unlock()
	broker.Rec("Subscribe -- %s: %s", client.id, topic)
	subs, found := broker.subscriptions[topic]
	if !found {
		broker.subscriptions[topic] = SubscriptionList{client}
	} else {
		for _, c := range subs {
			if c == client {
				return
			}
		}
		broker.subscriptions[topic] = append(subs, client)
	}

	// send all retained messages for this subscription
	for t, message := range broker.retained {
		if topicMatch(topic, t) {
			broker.Rec("(retain) -> %s: %s", message.Topic, FormatMQTTMessage(message))
			broker.queueMessage(client, message)
		}
	}

}

func (broker *FakeMQTTBroker) Unsubscribe(client *FakeMQTTClient, topic string) {
	broker.Lock()
	defer broker.Unlock()
	broker.Rec("Unsubscribe -- %s: %s", client.id, topic)
	subs, found := broker.subscriptions[topic]
	if !found {
		return
	} else {
		newSubs := make(SubscriptionList, 0, len(subs))
		for _, c := range subs {
			if c != client {
				newSubs = append(newSubs, c)
			}
		}
		broker.subscriptions[topic] = newSubs
	}
}

func (broker *FakeMQTTBroker) MakeClient(id string) (client *FakeMQTTClient) {
	client = &FakeMQTTClient{
		id:          id,
		started:     false,
		broker:      broker,
		callbackMap: make(map[string][]wbgong.MQTTMessageHandler),
		ready:       make(chan struct{}),
	}
	if broker.waitForRetained {
		broker.readyChannels = append(broker.readyChannels, client.ready)
	}

	broker.numClients += 1
	if broker.numClients == 1 {
		broker.Start()
	}

	return client
}

func (broker *FakeMQTTBroker) removeClient() {
	broker.numClients -= 1
	if broker.numClients == 0 {
		broker.Stop()
	}
}

type FakeMQTTClient struct {
	sync.Mutex
	id          string
	started     bool
	broker      *FakeMQTTBroker
	callbackMap map[string][]wbgong.MQTTMessageHandler
	ready       chan struct{}
}

func (client *FakeMQTTClient) receive(message wbgong.MQTTMessage) {
	client.Lock()

	// make a deep copy of callbacks map
	localMap := make(map[string][]wbgong.MQTTMessageHandler)
	for topic, handlers := range client.callbackMap {
		localHndlrs := make([]wbgong.MQTTMessageHandler, len(handlers))
		for i := range handlers {
			localHndlrs[i] = handlers[i]
		}
		localMap[topic] = localHndlrs
	}
	client.Unlock()

	for topic, handlers := range localMap {
		if !topicMatch(topic, message.Topic) {
			continue
		}
		for _, handler := range handlers {
			handler(message)
		}
	}
}

func (client *FakeMQTTClient) WaitForRetained(c func()) {
	if client.broker.waitForRetained {
		go func() {
			<-client.ready
			client.broker.PushRetainHack(c)
		}()
	} else {
		client.broker.PushRetainHack(c)
	}
}

func (client *FakeMQTTClient) Start() {
	if client.started {
		return
	}
	client.started = true
	if !client.broker.waitForRetained {
		close(client.ready)
	}
}

func (client *FakeMQTTClient) Stop() {
	client.ensureStarted()
	client.started = false
	client.broker.Rec("stop: %s", client.id)
	client.broker.removeClient()
}

func (client *FakeMQTTClient) ensureStarted() {
	if !client.started {
		log.Panicf("%s: client not started", client.id)
	}
}

func (client *FakeMQTTClient) Publish(message wbgong.MQTTMessage) {
	client.ensureStarted()
	client.broker.Publish(client.id, message)
}

func (client *FakeMQTTClient) Subscribe(callback wbgong.MQTTMessageHandler, topics ...string) {
	client.Lock()
	defer client.Unlock()
	client.ensureStarted()
	for _, topic := range topics {
		client.broker.Subscribe(client, topic)
		handlerList, found := client.callbackMap[topic]
		if found {
			client.callbackMap[topic] = append(handlerList, callback)
		} else {
			client.callbackMap[topic] = []wbgong.MQTTMessageHandler{callback}
		}
	}
}

func (client *FakeMQTTClient) Unsubscribe(topics ...string) {
	client.Lock()
	defer client.Unlock()
	client.ensureStarted()
	for _, topic := range topics {
		client.broker.Unsubscribe(client, topic)
		delete(client.callbackMap, topic)
	}
}

type FakeMQTTFixture struct {
	*Recorder
	Broker *FakeMQTTBroker
}

func NewFakeMQTTFixture(t *testing.T) *FakeMQTTFixture {
	rec := NewRecorder(t)
	return &FakeMQTTFixture{
		Recorder: rec,
		Broker:   NewFakeMQTTBroker(t, rec),
	}
}

// Shuffles slice of messages
// Sometimes it's quite useful to make MQTT related tests more representable
func ShuffleMessages(arr []wbgong.MQTTMessage) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := len(arr) - 1; i > 0; i-- {
		j := r.Intn(i)
		arr[i], arr[j] = arr[j], arr[i]
	}
}
