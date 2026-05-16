package di

import (
	"auth-micro-service/internal/rabbitMQ"
	"fmt"
)

func (d *DI) GetPublisher() *rabbitMQ.Publisher {
	publisher, err := rabbitMQ.New(d.GetRmq())

	if err != nil {
		d.logger.Error(fmt.Sprintf("failed to get publisher %v", err))
	}

	return publisher
}
