package config

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerRabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func InitRabbitMQ(cfg ConsumerConfig) ConsumerRabbitMQ {
	// create rabbitmq connection
	connection, err := amqp.Dial(cfg.BrokerURL)
	if err != nil {
		log.Fatalln("failed to connect to rabbitMQ broker, ", err)
	}

	// open rabbitMQ channel
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalln("failed to open rabbitMQ channel, ", err)
	}

	// declare exchange
	err = channel.ExchangeDeclare(cfg.ExchangeName, cfg.ExchangeType, true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare exchange, ", err)
	}

	// declare queue
	_, err = channel.QueueDeclare(cfg.QueueName, true, false, false, false, amqp.Table{"x-dead-letter-exchange": cfg.DLXName})
	if err != nil {
		log.Fatalln("failed to declare queue, ", err)
	}

	// bind queue to exchange
	err = channel.QueueBind(cfg.QueueName, cfg.QueueRoutingKey, cfg.ExchangeName, false, nil)
	if err != nil {
		log.Fatalln("failed to bind queue, ", err)
	}

	// set prefetch count => how many message will be published to Consumer
	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalln("failed to set prefetch count, ", err)
	}

	// configure dead letter exchange and queue for dropped message to be retried
	err = channel.ExchangeDeclare(cfg.DLXName, cfg.DLXType, true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare dead letter exchange, ", err)
	}

	_, err = channel.QueueDeclare(cfg.DLQName, true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare dead letter queue, ", err)
	}

	err = channel.QueueBind(cfg.DLQName, "", cfg.DLXName, false, nil)
	if err != nil {
		log.Fatalln("failed to bind dead letter queue to dead letter exchange, ", err)
	}

	return ConsumerRabbitMQ{
		Connection: connection,
		Channel:    channel,
	}
}
