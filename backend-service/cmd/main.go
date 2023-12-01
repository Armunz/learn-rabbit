package main

import (
	"learn-rabbit/backend-service/internal/config"
	"learn-rabbit/backend-service/internal/controller"
	"learn-rabbit/backend-service/internal/repository"
	"learn-rabbit/backend-service/internal/service"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("init config...")
	cfg := config.InitConfig()

	log.Println("init repository...")
	repo := repository.NewBackendRepository(cfg.ProducerURL, cfg.TopicName, cfg.RoutingKey, cfg.RepoTimeout)

	log.Println("init service...")
	s := service.NewBackendService(repo)

	log.Println("init controller...")
	c := controller.NewBackendController(s, cfg.APITimeout)

	app := fiber.New()
	c.Route(app)

	// Create a channel to listen for OS signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Run Fiber server in a separate goroutine
	go func() {
		if err := app.Listen(":8888"); err != nil {
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

}
