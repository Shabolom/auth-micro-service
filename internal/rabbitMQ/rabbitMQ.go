package rabbitMQ

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp.Channel
}

const (
	AuthExchange = "auth.events"
)

func New(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.New("failed to open a channel")
	}

	err = ch.ExchangeDeclare(
		AuthExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		return nil, errors.New("failed to declare an exchange")
	}

	return &Publisher{
		ch: ch,
	}, nil
}
