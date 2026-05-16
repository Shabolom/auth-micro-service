package rabbitMQ

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (p *Publisher) Publish(ctx context.Context, routingKey string, contentType string, body []byte) error {
	return p.ch.PublishWithContext(
		ctx,
		"auth.events",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
