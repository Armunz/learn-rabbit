package service

import (
	"context"
	"learn-rabbit/producer/internal/repository"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ProducerServiceImpl struct {
	repo repository.ProducerRepository
}

type ProducerService interface {
	Publish(ctx context.Context, exchangeName string, routingKey string, data []byte) error
}

func NewProducerService(repo repository.ProducerRepository) ProducerService {
	return &ProducerServiceImpl{
		repo: repo,
	}
}

// Publish implements ProducerService
func (s *ProducerServiceImpl) Publish(ctx context.Context, exchangeName string, routingKey string, data []byte) error {
	message := amqp.Publishing{
		DeliveryMode: 0, // transient (which means message is not persistent)
		Timestamp:    time.Now(),
		Body:         data,
	}

	return s.repo.Publish(ctx, exchangeName, routingKey, message)
}
