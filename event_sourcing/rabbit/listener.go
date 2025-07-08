package rabbit

import (
	"encoding/json"
	"github.com/IntelXLabs-LLC/go-common/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ListenerFunc[B any] func(body *B) error

type ListenerI[B any] interface {
	ListenAsync()
	Close()
	AssignHandlerFunction(fn ListenerFunc[B])
}

type Listener[B any] struct {
	ch       *amqp.Channel
	cl       *Client
	topic    string
	actionFn ListenerFunc[B]
	closer   chan struct{}
	body     *B
}

func (l *Listener[B]) AssignHandlerFunction(fn ListenerFunc[B]) {
	l.actionFn = fn
}

func (l *Listener[B]) ListenAsync() {
	consumerChannel, err := l.ch.Consume(l.topic, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		for msg := range consumerChannel {
			// Call the action function with the message body
			err := json.Unmarshal(msg.Body, &l.body)
			if err != nil {
				logger.Error(err, "Cannot Unmarshal inside Listener")
			}

			err = l.actionFn(l.body)
			if err != nil {
				logger.Error(err, "Cannot Run Action function inside Listener")
			}
		}
	}()

	<-l.closer
}

func (l *Listener[B]) Close() {
	l.closer <- struct{}{}
}
