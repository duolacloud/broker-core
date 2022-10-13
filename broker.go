package broker

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Init(...Option) error
	Options() Options
	Address() string
	Connect() error
	Disconnect() error
	Publish(topic string, m *Message, opts ...PublishOption) error
	Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	String() string
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(Event) error

// Message is a message send/received from the broker.
type Message struct {
	Header map[string]string
	Body   []byte
}

// Event is given to a subscription handler for processing.
type Event interface {
	Topic() string
	Message() *Message
	Ack() error
	Error() error
}

// Subscriber is a convenience return type for the Subscribe method.
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}

// Marshaler is a simple encoding interface used for the broker/transport
// where headers are not supported by the underlying implementation.
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}
