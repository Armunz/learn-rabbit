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

	//
	args := amqp.Table{
		"x-dead-letter-exchange":    cfg.DLXName,
		"x-dead-letter-routing-key": cfg.DLXRoutingKey,
	}
	// declare queue
	_, err = channel.QueueDeclare(cfg.QueueName, true, false, false, false, args)
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
	err = channel.ExchangeDeclare(cfg.DLXName, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare dead letter exchange, ", err)
	}

	_, err = channel.QueueDeclare(cfg.DLQName, true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare dead letter queue, ", err)
	}

	// bind dlx to dlq
	err = channel.QueueBind(cfg.DLQName, cfg.DLXRoutingKey, cfg.DLXName, false, nil)
	if err != nil {
		log.Fatalln("failed to bind dead letter queue to dead letter exchange, ", err)
	}

	// // Set the main queue to use the dead letter exchange
	// args := amqp.Table{
	// 	"x-dead-letter-exchange":    cfg.DLXName,       // Dead letter exchange name
	// 	"x-dead-letter-routing-key": cfg.DLXRoutingKey, // Dead letter routing key
	// }

	// // Bind the dead lettering arguments to the main queue
	// err = channel.QueueBind(
	// 	cfg.QueueName,       // Queue name
	// 	cfg.QueueRoutingKey, // Routing key
	// 	cfg.ExchangeName,    // Exchange name
	// 	false,               // No-wait
	// 	args,                // Arguments with dead lettering configuration
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to bind the dead lettering arguments to the main queue: %s", err)
	// }

	return ConsumerRabbitMQ{
		Connection: connection,
		Channel:    channel,
	}
}
