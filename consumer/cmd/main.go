package main

import (
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

	// init repository
	log.Println("Init repository...")
	repo := repository.NewUserRepository(dbUser, cfg.MYSQLQueryTimeoutMs)

	// init service
	log.Println("Init service...")
	consumerService := service.NewConsumerService(repo)

	// init controller
	log.Println("Init controller...")
	consumerController := controller.NewConsumerController(consumerService)

	// init consumer rabbitMQ
	log.Println("Init consumer rabbitMQ...")
	var consumerRabbitMQ *config.ConsumerRabbitMQ
	if cfg.IsUsingCluster {
		consumerRabbitMQ = config.InitRabbitMQCluster(cfg, consumerController)
	} else {
		consumerRabbitMQ = config.InitRabbitMQ(cfg)
	}

	log.Println("Start consuming message...")

	// Create a channel to communicate termination signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Handle termination signal in a separate goroutine
	go func() {
		<-stopChan
		log.Println("Received termination signal. Gracefully stopping consumer...")

		consumerRabbitMQ.StopConsuming()

		log.Printf("Stop %s...\n", cfg.ConsumerName)
		// Stop consuming messages
		if err := consumerRabbitMQ.Channel.Cancel(cfg.ConsumerName, false); err != nil {
			log.Fatalf("Error cancelling consumer: %s", err)
		}

		log.Printf("Stop %s...\n", cfg.ConsumerDroppedMessageName)
		if err := consumerRabbitMQ.Channel.Cancel(cfg.ConsumerDroppedMessageName, false); err != nil {
			log.Fatalf("Error cancelling dropped message consumer: %s", err)
		}

		log.Println("Close consumer channel...")
		if err := consumerRabbitMQ.Channel.Close(); err != nil {
			log.Fatalln("Failed to close consumer rabbitMQ channel ,", err)
		}

		log.Println("Close consumer connection...")
		if err := consumerRabbitMQ.Connection.Close(); err != nil {
			log.Fatalln("Failed to close consumer rabbitMQ connection ,", err)
		}

		log.Println("Close MYSQL database connection...")
		config.CleanUp(dbUser)

		// Exit the program
		os.Exit(0)
	}()

	// start consume message
	consumerRabbitMQ.StartConsume(cfg, consumerController)
	// go func(c *config.ConsumerRabbitMQ) {
	// 	messages, err := c.Channel.Consume(cfg.QueueName, cfg.ConsumerName, false, false, false, false, nil)
	// 	if err != nil {
	// 		log.Fatalln("failed to consume message from producer, ", err)
	// 	}

	// 	for delivery := range messages {
	// 		consumerController.Handle(context.Background(), delivery)
	// 	}
	// }(consumerRabbitMQ)

	// // start consume dropped message
	// go func(c *config.ConsumerRabbitMQ) {
	// 	droppedMessages, err := c.Channel.Consume(cfg.DLQName, cfg.ConsumerDroppedMessageName, false, false, false, false, nil)
	// 	if err != nil {
	// 		log.Fatalln("failed to consume dropped message, ", err)
	// 	}

	// 	for droppedDelivery := range droppedMessages {
	// 		consumerController.HandleDroppedMessage(context.Background(), droppedDelivery)
	// 	}
	// }(consumerRabbitMQ)

	// Keep the main goroutine running
	select {}
}
