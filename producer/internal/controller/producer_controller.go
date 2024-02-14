package controller

import (
	"context"
	"encoding/json"
	"learn-rabbit/producer/internal/service"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type producerControllerImpl struct {
	producerService service.ProducerService
	timeoutMs       int
}

func NewProducerController(producerService service.ProducerService, timeoutMs int) *producerControllerImpl {
	return &producerControllerImpl{
		producerService: producerService,
		timeoutMs:       timeoutMs,
	}
}

func (controller *producerControllerImpl) Route(app *fiber.App) {
	app.Post("/user", controller.Handler)
}

func (controller *producerControllerImpl) Handler(c *fiber.Ctx) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(controller.timeoutMs)*time.Millisecond)
	defer cancel()

	var request Request
	err := c.BodyParser(&request)
	if err != nil {
		log.Println("failed to parse request body, ", err)
		return err
	}

	data, err := json.Marshal(request.Data)
	if err != nil {
		log.Println("failed to marshal request data, ", err)
		return err
	}

	err = controller.producerService.Publish(ctxTimeout, request.ExchangeName, request.RoutingKey, data)
	if err != nil {
		return err
	}

	c.Status(fiber.StatusCreated)
	return nil
}
