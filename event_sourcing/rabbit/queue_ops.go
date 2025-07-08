package rabbit

import (
	"errors"
	"github.com/IntelXLabs-LLC/go-common/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Client) createQueue(topic string) (*amqp.Channel, error) {
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func NewPublisher(c *Client, topic string) (PublisherI, error) {
	// Declare new Publisher
	p := &Publisher{
		ch: nil,
		cl: c,
	}

	confirms := make(chan amqp.Confirmation)

	ch, err := p.cl.createQueue(topic)
	if err != nil {
		panic(err)
	}
	p.ch = ch

	p.ch.NotifyPublish(confirms)

	go func() {
		for confirm := range confirms {
			if confirm.Ack {
				logger.Info("Message published successfully")
			} else {
				logger.Error(errors.New("failed to publish message"))
			}
		}
	}()

	err = p.ch.Confirm(false)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func NewListener[B any](c *Client, topic string) (ListenerI[B], error) {
	l := &Listener[B]{
		ch:       nil,
		cl:       c,
		topic:    topic,
		actionFn: nil,
		closer:   make(chan struct{}),
	}

	ch, err := l.cl.createQueue(l.topic)
	if err != nil {
		panic(err)
	}
	// Assign the channel for the first time
	l.ch = ch

	return l, nil
}
