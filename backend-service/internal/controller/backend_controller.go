package controller

import (
	"context"
	"learn-rabbit/backend-service/internal/model"
	"learn-rabbit/backend-service/internal/service"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type backendController struct {
	backendService service.BackendService
	timeoutMs      int
}

func NewBackendController(backendService service.BackendService, timeoutMs int) *backendController {
	return &backendController{
		backendService: backendService,
		timeoutMs:      timeoutMs,
	}
}

func (b *backendController) SaveUser(c *fiber.Ctx) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(b.timeoutMs)*time.Millisecond)
	defer cancel()

	userRequest := new(model.UserRequest)
	if err := c.BodyParser(userRequest); err != nil {
		log.Println("failed to parse request, ", err)
		return err
	}

	if err := b.backendService.SaveUser(ctxTimeout, *userRequest); err != nil {
		log.Println("failed to save user, ", err)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(model.GeneralResponse{
		Code:    201,
		Message: "User Created",
	})
}

func (b *backendController) Route(app *fiber.App) {
	app.Post("/user", b.SaveUser)
}
