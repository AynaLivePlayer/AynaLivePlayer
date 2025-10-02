package eventbus

type Event struct {
	// Id is the event id
	Id string
	// Channel if channel is empty, then event is broadcast
	Channel string
	// EchoId is used for callback, if echo is not empty
	// the caller is expecting a callback
	EchoId string
	// Data any data struct
	Data interface{}
}

// HandlerFunc event handler, should be non-blocking
type HandlerFunc func(event *Event)

// Subscriber is client to the bus
type Subscriber interface {
	// Subscribe will run this handler asynchronous when an event received;
	// event will still come sequentially for each handler. which means before previous
	// event has finished, the same handler should not be called.
	// if channel is not empty, the handler will not receive event from other channel, however,
	// broadcast event (channel is empty) will still be passed to the handler
	Subscribe(channel string, eventId string, handlerName string, fn HandlerFunc) error
	// SubscribeAny is Subscribe with empty channel. this function will subscribe to event from any channel
	SubscribeAny(eventId string, handlerName string, fn HandlerFunc) error
	// SubscribeOnce will run handler once, and delete handler internally
	SubscribeOnce(channel string, eventId string, handlerName string, fn HandlerFunc) error
	// Unsubscribe just remove handler for the bus
	Unsubscribe(eventId string, handlerName string) error
}

type Publisher interface {
	// Publish basically a wrapper to PublishEvent
	Publish(eventId string, data interface{}) error
	// PublishToChannel publish event to a specific channel, basically another wrapper to PublishEvent
	PublishToChannel(channel string, eventId string, data interface{}) error
	// PublishEvent publish an event
	PublishEvent(event *Event) error
}

// Caller is special usage of a Publisher
type Caller interface {
	Call(pubEvtId string, data interface{}, subEvtId string) (*Event, error)
	Reply(req *Event, eventId string, data interface{}) error
}

type Controller interface {
	// Start will start to push events to subscribers,
	// Publisher should be able to publish events before bus started
	Start() error
	// Wait will wait all event to be executed
	Wait() error
	// Stop will stop controller immediately
	Stop() error
}

type Bus interface {
	Controller
	Publisher
	Subscriber
	Caller
}
