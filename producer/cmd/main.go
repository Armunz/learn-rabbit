package main

import (
	"learn-rabbit/producer/internal/config"
	"learn-rabbit/producer/internal/controller"
	"learn-rabbit/producer/internal/repository"
	"learn-rabbit/producer/internal/service"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// init config
	log.Println("Init config...")
	cfg := config.InitConfig()

	// init rabbitMQ
	log.Println("Init producer rabbitMQ...")
	producerRabbitMQ := config.InitRabbitMQ(cfg)

	// init repository
	log.Println("Init repository...")
	repo := repository.NewProducerRepository(producerRabbitMQ.Channel, cfg.RepoTimeoutMs)

	// init service
	log.Println("Init service...")
	producerService := service.NewProducerService(repo)

	// init controller
	log.Println("Init controller...")
	producerController := controller.NewProducerController(producerService, cfg.APITimeoutMs)

	// init fiber
	app := fiber.New()

	// setup routing
	producerController.Route(app)

	// Create a channel to listen for OS signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Run Fiber server in a separate goroutine
	go func() {
		if err := app.Listen(":9999"); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}()

	// Wait for a termination signal
	<-stopChan
	log.Println("Received termination signal. Shutting down server...")

	// Shutdown server
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %s", err)
	}

	log.Println("Server gracefully stopped.")

	if err := producerRabbitMQ.Channel.Close(); err != nil {
		log.Fatalln("Failed to close rabbitMQ channel, ", err)
	}

	if err := producerRabbitMQ.Connection.Close(); err != nil {
		log.Fatalln("Failed to close rabbitMQ connection, ", err)
	}

	log.Println("RabbitMQ connection and channel closed.")

	// Exit the application
	os.Exit(0)
}
