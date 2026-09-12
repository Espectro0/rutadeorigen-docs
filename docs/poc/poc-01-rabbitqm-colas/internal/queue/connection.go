package queue

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func Connect(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to connnect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Failed to open a channel: %v", err)
	}

	return &Connection{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (c *Connection) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}

	if c.Conn != nil {
		c.Conn.Close()
	}
}

func (c *Connection) DeclareExchange(name, kind string) error {
	return c.Channel.ExchangeDeclare(
		name,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *Connection) DeclareQueue(queueName string, args amqp.Table) (amqp.Queue, error) {
	return c.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		args,
	)
}

func (c *Connection) BindQueue(queueName, routingKey, exchangeName string) error {
	return c.Channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
}

func (c *Connection) Publish(ctx context.Context, exchange, routingKey string, headers amqp.Table, body []byte) error {
	return c.Channel.PublishWithContext(ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Headers:      headers,
			Body:         body,
		},
	)
}
