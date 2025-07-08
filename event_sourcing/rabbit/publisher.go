package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

var once sync.Once

type PublisherI interface {
	Publish(t TopicI) error
	Close() error
}

type Publisher struct {
	ch *amqp.Channel
	cl *Client
}

// Publish publishes the data to the queue
func (p *Publisher) Publish(t TopicI) error {
	// Publish with Context
	return p.ch.Publish("", t.GetTopicName(), false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         t.GetBody(),
	})
}

func (p *Publisher) Close() error {
	if p.ch != nil {
		return p.ch.Close()
	}
	return nil
}
