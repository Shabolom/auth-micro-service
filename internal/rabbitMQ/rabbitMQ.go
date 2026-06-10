package rabbitMQ

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp.Channel
}

const (
	TEXTTYPE string = "text/plain"
	JSONTYPE string = "application/json"
)

func New(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, errors.New("failed to open a channel")
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
		_ = conn.Close()
		return nil, errors.New("failed to declare an exchange")
	}

	_, err = ch.QueueDeclare(
		"auth.register",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, errors.New("failed to declare a queue")
	}

	_, err = ch.QueueDeclare(
		"auth.login.logs",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, errors.New("failed to declare a queue")
	}

	err = ch.QueueBind(
		"auth.register",
		"register",
		"auth.events",
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, errors.New("failed to bind a queue")
	}

	err = ch.QueueBind(
		"auth.login.logs",
		"login",
		"auth.events",
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, errors.New("failed to bind a queue")
	}

	return &Publisher{
		ch: ch,
	}, nil
}
