package controller

import (
	"context"
	"encoding/json"
	"learn-rabbit/producer/internal/service"
	"time"

	"github.com/gofiber/fiber"
)

type ProducerControllerImpl struct {
	producerService service.ProducerService
	timeoutMs       int
}

func NewProducerController(producerService service.ProducerService, timeoutMs int) *ProducerControllerImpl {
	return &ProducerControllerImpl{
		producerService: producerService,
		timeoutMs:       timeoutMs,
	}
}

func (controller *ProducerControllerImpl) Route(app *fiber.App) {
	app.Post("/user", controller.Handler)
}

func (controller *ProducerControllerImpl) Handler(c *fiber.Ctx) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(controller.timeoutMs)*time.Millisecond)
	defer cancel()

	request := new(Request)
	err := c.BodyParser(request)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).Send(err)
	}

	data, err := json.Marshal(request.Data)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).Send(err)
	}

	err = controller.producerService.Publish(ctxTimeout, request.ExchangeName, request.RoutingKey, data)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).Send(err)
	}

	c.Status(fiber.StatusOK)
}
