package repository

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ProducerRepositoryImpl struct {
	rabbitMQChannel *amqp.Channel
	timeoutMs       int
}

type ProducerRepository interface {
	Publish(ctx context.Context, exchangeName string, routingKey string, message amqp.Publishing) error
}

func NewProducerRepository(rabbitMQChannel *amqp.Channel, timeoutMs int) ProducerRepository {
	return &ProducerRepositoryImpl{
		rabbitMQChannel: rabbitMQChannel,
		timeoutMs:       timeoutMs,
	}
}

// Publish implements ProducerRepository
func (r *ProducerRepositoryImpl) Publish(ctx context.Context, exchangeName string, routingKey string, message amqp.Publishing) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	return r.rabbitMQChannel.PublishWithContext(ctxTimeout, exchangeName, routingKey, false, false, message)
}
