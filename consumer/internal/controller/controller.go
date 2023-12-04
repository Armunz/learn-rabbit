package controller

import (
	"context"
	"learn-rabbit/consumer/internal/service"
	"log"
	"time"

	"github.com/avast/retry-go"

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
		return
	}

	// when there is no error, then acknowledge the message
	if err := delivery.Ack(false); err != nil {
		log.Println("failed to acknowledge the message, ", err)
	}

	log.Println("message acknowledged")

	return
}

func (c *ConsumerControllerImpl) HandleDroppedMessage(ctx context.Context, delivery amqp.Delivery) {
	log.Println("acknowledging dropped message...")

	body := delivery.Body
	log.Println("Dropped Message: ", string(body))

	// we can implement retry mechanism when dead message failed to handle
	err := retry.Do(
		func() error {
			return c.consumerService.Handle(ctx, body)
		},
		retry.Attempts(3),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retrying request after error: %v", err)
		}),
		retry.Delay(time.Duration(1000)*time.Millisecond), // set default delay
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return retry.BackOffDelay(n, err, config) // increase the delay time between consecutive retries
		}),
	)
	if err != nil {
		// if max attempt still fails, then we can log the failed message to db
		log.Println("Failed message: ", string(body))
	}

	if err := delivery.Ack(false); err != nil {
		log.Println("failed to acknowledge the dropped message, ", err)
	}

	log.Println("dropped message acknowledged")
}
