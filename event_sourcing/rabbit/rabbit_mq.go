package rabbit

import (
	"fmt"
	"github.com/IntelXLabs-LLC/go-common/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/url"
)

type Client struct {
	Connection *amqp.Connection
}

func NewClient(cfg *config.RabbitMQConfig) (*Client, error) {
	c := new(Client)
	err := c.Connect(cfg)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Connect(cfg *config.RabbitMQConfig) error {
	if cfg == nil {
		return ErrInvalidConfig
	}
	uri := fmt.Sprintf("amqps://%s:%s@%s:%s", url.PathEscape(cfg.Username), url.PathEscape(cfg.Password), cfg.Host, url.PathEscape(cfg.Port))
	// new code to encode the uri
	var err error
	c.Connection, err = amqp.Dial(uri)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	if c.Connection != nil {
		return c.Connection.Close()
	}
	return nil
}
