package rabbitMQ

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp.Channel
}

const (
	TEXTTYPE = "text/plain"
	JSONTYPE = "application/json"
)

func New(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		"auth.events",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to declare exchange auth.events: %w", err)
	}

	return &Publisher{ch: ch}, nil
}
