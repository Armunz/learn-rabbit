package repository

import (
	"context"
	"learn-rabbit/producer/internal/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type producerRepositoryImpl struct {
	rabbitMQProducer *config.ProducerRabbitMQ
	timeoutMs        int
}

type ProducerRepository interface {
	Publish(ctx context.Context, exchangeName string, routingKey string, message amqp.Publishing) error
}

func NewProducerRepository(rabbitMQProducer *config.ProducerRabbitMQ, timeoutMs int) ProducerRepository {
	return &producerRepositoryImpl{
		rabbitMQProducer: rabbitMQProducer,
		timeoutMs:        timeoutMs,
	}
}

// Publish implements ProducerRepository
func (r *producerRepositoryImpl) Publish(ctx context.Context, exchangeName string, routingKey string, message amqp.Publishing) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	return r.rabbitMQProducer.Channel.PublishWithContext(ctxTimeout, exchangeName, routingKey, false, false, message)
}
