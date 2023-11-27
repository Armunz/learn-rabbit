package cmd

import (
	"learn-rabbit/producer/internal/config"
	"learn-rabbit/producer/internal/controller"
	"learn-rabbit/producer/internal/repository"
	"learn-rabbit/producer/internal/service"

	"github.com/gofiber/fiber"
)

func main() {
	// init config
	cfg := config.InitConfig()

	// init rabbitMQ
	rbmqChannel := config.InitRabbitMQ(cfg)

	// init repository
	repo := repository.NewProducerRepository(rbmqChannel, cfg.RepoTimeoutMs)

	// init service
	producerService := service.NewProducerService(repo)

	// init controller
	producerController := controller.NewProducerController(producerService, cfg.APITimeoutMs)

	// init fiber
	app := fiber.New()

	// setup routing
	producerController.Route(app)

	err := app.Listen("9999")
	if err != nil {
		panic(err)
	}
}
