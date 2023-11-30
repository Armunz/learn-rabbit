package main

import (
	"context"
	"learn-rabbit/consumer/internal/config"
	"learn-rabbit/consumer/internal/controller"
	"learn-rabbit/consumer/internal/repository"
	"learn-rabbit/consumer/internal/service"
	"log"
	"os"
	"os/signal"
)

func main() {
	// init config
	log.Println("Init config...")
	cfg := config.InitConfig()

	// init mysql
	log.Println("Init mysql database...")
	dbUser := config.NewMySQLDatabase(cfg)

	// init consumer rabbitMQ
	log.Println("Init consumer rabbitMQ...")
	consumerRabbitMQ := config.InitRabbitMQ(cfg)

	// init repository
	log.Println("Init repository...")
	repo := repository.NewUserRepository(dbUser, cfg.MYSQLQueryTimeoutMs)

	// init service
	log.Println("Init service...")
	consumerService := service.NewConsumerService(repo)

	// init controller
	consumerController := controller.NewConsumerController(consumerService)

	messages, err := consumerRabbitMQ.Channel.Consume(cfg.QueueName, cfg.ConsumerName, false, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to consume message from producer, ", err)
	}

	droppedMessages, err := consumerRabbitMQ.Channel.Consume(cfg.DLQName, cfg.ConsumerName, false, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to consume dropped message, ", err)
	}

	// Create a channel to communicate termination signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Handle termination signal in a separate goroutine
	go func() {
		<-stopChan
		log.Println("Received termination signal. Gracefully stopping consumer...")

		// Stop consuming messages
		if err := consumerRabbitMQ.Channel.Cancel(cfg.ConsumerName, false); err != nil {
			log.Printf("Error cancelling consumer: %s", err)
		}
	}()

	// start consume message
	go func() {
		for delivery := range messages {
			consumerController.Handle(context.Background(), delivery)
		}
	}()

	// start consume dropped message
	go func() {
		for delivery := range droppedMessages {
			consumerController.HandleDroppedMessage(context.Background(), delivery)
		}
	}()

	// Keep the main goroutine running
	select {}
}
