package main

import (
	"learn-rabbit/backend-service/internal/config"
	"learn-rabbit/backend-service/internal/controller"
	"learn-rabbit/backend-service/internal/repository"
	"learn-rabbit/backend-service/internal/service"
	"log"

	"github.com/gofiber/fiber"
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
	app.Listen("9999")
}
