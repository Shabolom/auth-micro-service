package di

import (
	"auth-micro-service/internal/rabbitMQ"

	"go.uber.org/zap"
)

func (d *DI) GetPublisher() *rabbitMQ.Publisher {
	publisher, err := rabbitMQ.New(d.GetRmq())

	if err != nil {
		d.logger.Fatal("failed to get publisher", zap.Error(err))
	}

	return publisher
}
