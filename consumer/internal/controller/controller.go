package controller

import (
	"context"
	"learn-rabbit/consumer/internal/service"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerControllerImpl struct {
	consumerService service.ConsumerService
}

type ConsumerController interface {
	Handle(ctx context.Context, delivery amqp.Delivery)
	HandleDroppedMessage(ctx context.Context, delivery amqp.Delivery)
}

func NewConsumerController(consumerService service.ConsumerService) ConsumerController {
	return &ConsumerControllerImpl{
		consumerService: consumerService,
	}
}

// Handle implements ConsumerController.
func (c *ConsumerControllerImpl) Handle(ctx context.Context, delivery amqp.Delivery) {
	log.Println("acknowledging message...")

	body := delivery.Body
	if err := c.consumerService.Handle(ctx, body); err != nil {
		// when error occurs, then throw message to dead letter queue
		if err := delivery.Nack(false, false); err != nil {
			log.Println("failed to neglect the message, ", err)
		}
	}

	// when there is no error, then acknowledge the message
	if err := delivery.Ack(false); err != nil {
		log.Println("failed to acknowledge the message, ", err)
	}

	return
}

func (c *ConsumerControllerImpl) HandleDroppedMessage(ctx context.Context, delivery amqp.Delivery) {
	log.Println("acknowledging dropped message...")

	body := delivery.Body
	log.Println("Dropped Message: ", string(body))

	if err := delivery.Ack(false); err != nil {
		log.Println("failed to acknowledge the dropped message, ", err)
	}
}
