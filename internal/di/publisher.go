package di

import (
	"auth-micro-service/internal/rebbitMQ"
	"fmt"
)

func (d *DI) GetPublisher() *rebbitMQ.Publisher {
	publisher, err := rebbitMQ.New(d.GetRmq())

	if err != nil {
		d.logger.Error(fmt.Sprintf("failed to get publisher %v", err))
	}

	return publisher
}
