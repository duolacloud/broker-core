package noop

import (
	"errors"
	"sync"

	"github.com/duolacloud/broker-core"
	"github.com/duolacloud/broker-core/marshaler"
	"github.com/google/uuid"
	"github.com/huandu/go-clone"
)

type H struct {
	handler    broker.Handler
	resultType any
}

type S struct {
	opts *broker.SubscribeOptions
	t    string
}

func (s *S) Options() broker.SubscribeOptions {
	return *s.opts
}

func (s *S) Topic() string {
	return s.t
}

func (s *S) Unsubscribe() error {
	return nil
}

type E struct {
	t string
	m any
	e error
}

func (e *E) Topic() string {
	return e.t
}

func (e *E) Message() any {
	return e.m
}

func (e *E) Ack() error {
	return nil
}

func (e *E) Error() error {
	return e.e
}

type NoopBroker struct {
	opts     *broker.Options
	handlers map[string]map[string][]H
	mutex    sync.Mutex
}

func NewBroker(opts ...broker.Option) broker.Broker {
	b := &NoopBroker{
		opts:     &broker.Options{},
		handlers: make(map[string]map[string][]H),
	}
	_ = b.Init(opts...)
	return b
}

func (b *NoopBroker) Init(opts ...broker.Option) error {
	for _, opt := range opts {
		opt(b.opts)
	}
	if b.opts.Codec == nil {
		b.opts.Codec = &marshaler.JsonMarshaler{}
	}
	return nil
}

func (b *NoopBroker) Options() broker.Options {
	return *b.opts
}

func (b *NoopBroker) Address() string {
	return ""
}

func (b *NoopBroker) Connect() error {
	return nil
}

func (b *NoopBroker) Disconnect() error {
	return nil
}

func (b *NoopBroker) Publish(topic string, msg any, opts ...broker.PublishOption) (err error) {
	eh := b.opts.ErrorHandler
	qhs, ok := b.handlers[topic]
	if !ok {
		return errors.New("topic not subscribed")
	}

	body, ok := msg.([]byte)
	if !ok {
		if body, err = b.opts.Codec.Marshal(msg); err != nil {
			return
		}
	}

	for _, hs := range qhs {
		if len(hs) == 0 {
			continue
		}
		h := hs[0]
		handler := h.handler
		resultType := h.resultType
		evt := &E{t: topic}

		if resultType != nil {
			evt.m = clone.Clone(resultType)
			err = b.opts.Codec.Unmarshal(body, evt.m)
			if err != nil && eh != nil {
				evt.m = body
				evt.e = err
				_ = eh(evt)
				continue
			}
		}
		if evt.m == nil {
			evt.m = body
		}
		_ = handler(evt)
	}
	return nil
}

func (b *NoopBroker) Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	subopts := &broker.SubscribeOptions{
		Queue: uuid.New().String(),
	}
	for _, opt := range opts {
		opt(subopts)
	}
	queue := subopts.Queue
	resultType := subopts.ResultType

	s := &S{t: topic, opts: subopts}

	qhs, ok := b.handlers[topic]
	if !ok {
		qhs = make(map[string][]H)
		b.handlers[topic] = qhs
	}

	_, ok = qhs[queue]
	if !ok {
		qhs[queue] = make([]H, 0)
	}
	qhs[queue] = append(qhs[queue], H{
		handler:    h,
		resultType: resultType,
	})
	b.handlers[topic] = qhs
	return s, nil
}

func (b *NoopBroker) String() string {
	return "noop"
}
