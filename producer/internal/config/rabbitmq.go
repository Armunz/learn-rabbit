package config

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ(cfg ProducerConfig) *amqp.Channel {
	// create rabbitmq connection
	connection, err := amqp.Dial(cfg.BrokerURL)
	if err != nil {
		log.Fatalln("failed to connect to rabbitMQ broker, ", err)
	}
	defer connection.Close()

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

	return channel
}
