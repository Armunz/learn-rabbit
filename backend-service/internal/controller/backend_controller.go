package controller

import (
	"context"
	"learn-rabbit/backend-service/internal/model"
	"learn-rabbit/backend-service/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BackendController struct {
	backendService service.BackendService
	timeoutMs      int
}

func NewBackendController(backendService service.BackendService, timeoutMs int) *BackendController {
	return &BackendController{
		backendService: backendService,
		timeoutMs:      timeoutMs,
	}
}

func (b *BackendController) SaveUser(c *fiber.Ctx) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(b.timeoutMs)*time.Millisecond)
	defer cancel()

	userRequest := new(model.UserRequest)
	if err := c.BodyParser(userRequest); err != nil {
		return err
	}

	if err := b.backendService.SaveUser(ctxTimeout, *userRequest); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(model.GeneralResponse{
		Code:    201,
		Message: "User Created",
	})
}

func (b *BackendController) Route(app *fiber.App) {
	app.Post("/user", b.SaveUser)
}
