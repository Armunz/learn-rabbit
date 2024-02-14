package service

import (
	"context"
	"learn-rabbit/producer/internal/repository"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type producerServiceImpl struct {
	repo repository.ProducerRepository
}

type ProducerService interface {
	Publish(ctx context.Context, exchangeName string, routingKey string, data []byte) error
}

func NewProducerService(repo repository.ProducerRepository) ProducerService {
	return &producerServiceImpl{
		repo: repo,
	}
}

// Publish implements ProducerService
func (s *producerServiceImpl) Publish(ctx context.Context, exchangeName string, routingKey string, data []byte) error {
	message := amqp.Publishing{
		DeliveryMode: 0, // transient (which means message is not persistent)
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         data,
	}

	return s.repo.Publish(ctx, exchangeName, routingKey, message)
}
