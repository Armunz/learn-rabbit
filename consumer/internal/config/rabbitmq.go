package config

import (
	"context"
	"learn-rabbit/consumer/internal/controller"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerRabbitMQ struct {
	Connection    *amqp.Connection
	Channel       *amqp.Channel
	mu            sync.Mutex
	stopConsuming chan struct{}
}

func InitRabbitMQ(cfg ConsumerConfig) *ConsumerRabbitMQ {
	// create rabbitmq connection
	connection, err := amqp.Dial(cfg.BrokerURL)
	if err != nil {
		log.Fatalln("failed to connect to rabbitMQ broker, ", err)
	}

	// open rabbitMQ channel
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalln("failed to open rabbitMQ channel, ", err)
	}

	// declare exchange
	err = channel.ExchangeDeclare(cfg.ExchangeName, cfg.ExchangeType, true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare exchange, ", err)
	}

	// argument for dead letter queue
	args := amqp.Table{
		"x-dead-letter-exchange":    cfg.DLXName,
		"x-dead-letter-routing-key": cfg.DLXRoutingKey,
	}
	// declare queue
	_, err = channel.QueueDeclare(cfg.QueueName, true, false, false, false, args)
	if err != nil {
		log.Fatalln("failed to declare queue, ", err)
	}

	// bind queue to exchange
	err = channel.QueueBind(cfg.QueueName, cfg.QueueRoutingKey, cfg.ExchangeName, false, nil)
	if err != nil {
		log.Fatalln("failed to bind queue, ", err)
	}

	// set prefetch count => how many message will be consumed by Consumer
	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalln("failed to set prefetch count, ", err)
	}

	// configure dead letter exchange and queue for dropped message to be retried
	err = channel.ExchangeDeclare(cfg.DLXName, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare dead letter exchange, ", err)
	}

	_, err = channel.QueueDeclare(cfg.DLQName, true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare dead letter queue, ", err)
	}

	// bind dlx to dlq
	err = channel.QueueBind(cfg.DLQName, cfg.DLXRoutingKey, cfg.DLXName, false, nil)
	if err != nil {
		log.Fatalln("failed to bind dead letter queue to dead letter exchange, ", err)
	}

	return &ConsumerRabbitMQ{
		Connection: connection,
		Channel:    channel,
	}
}

func InitRabbitMQCluster(cfg ConsumerConfig, consumerController controller.ConsumerController) *ConsumerRabbitMQ {
	rabbitMQCluster, err := DialCluster(cfg, consumerController)
	if err != nil {
		log.Fatalln("failed to connect rabbitMQ cluster, ", err)
	}
	// rabbitMQCluster.OpenChannel(cfg, consumerController)

	// declare exchange
	err = rabbitMQCluster.Channel.ExchangeDeclare(cfg.ExchangeName, cfg.ExchangeType, true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare exchange, ", err)
	}

	// argument for dead letter queue
	args := amqp.Table{
		"x-dead-letter-exchange":    cfg.DLXName,
		"x-dead-letter-routing-key": cfg.DLXRoutingKey,
		"x-dead-letter-strategy":    "at-least-once",                 // this is for dead letter message
		amqp.QueueTypeArg:           amqp.QueueTypeQuorum,            // queue type is quorum for cluster
		amqp.QueueOverflowArg:       amqp.QueueOverflowRejectPublish, // this is for dead letter message
	}
	// declare queue
	_, err = rabbitMQCluster.Channel.QueueDeclare(cfg.QueueName, true, false, false, false, args)
	if err != nil {
		log.Fatalln("failed to declare queue, ", err)
	}

	// bind queue to exchange
	err = rabbitMQCluster.Channel.QueueBind(cfg.QueueName, cfg.QueueRoutingKey, cfg.ExchangeName, false, nil)
	if err != nil {
		log.Fatalln("failed to bind queue, ", err)
	}

	// set prefetch count => how many message will be consumed by Consumer
	err = rabbitMQCluster.Channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalln("failed to set prefetch count, ", err)
	}

	// configure dead letter exchange and queue for dropped message to be retried
	err = rabbitMQCluster.Channel.ExchangeDeclare(cfg.DLXName, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare dead letter exchange, ", err)
	}

	// declare dead letter queue with type quorum also
	_, err = rabbitMQCluster.Channel.QueueDeclare(cfg.DLQName, true, false, false, false, amqp.Table{
		amqp.QueueTypeArg: amqp.QueueTypeQuorum,
	})
	if err != nil {
		log.Fatalln("failed to declare dead letter queue, ", err)
	}

	// bind dlx to dlq
	err = rabbitMQCluster.Channel.QueueBind(cfg.DLQName, cfg.DLXRoutingKey, cfg.DLXName, false, nil)
	if err != nil {
		log.Fatalln("failed to bind dead letter queue to dead letter exchange, ", err)
	}

	return rabbitMQCluster
}

// based on https://github.com/sirius1024/go-amqp-reconnect/tree/master
func DialCluster(cfg ConsumerConfig, consumerController controller.ConsumerController) (*ConsumerRabbitMQ, error) {
	nodeSequence := 0
	var conn *amqp.Connection
	var err error
	for {
		conn, err = amqp.Dial(cfg.BrokerClusterURLs[nodeSequence])
		if err == nil {
			break
		}
		newSeq := next(cfg.BrokerClusterURLs, nodeSequence)
		nodeSequence = newSeq
		time.Sleep(time.Duration(1) * time.Second)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalln("failed to open rabbitMQ cluster channel, ", err)
	}

	rabbitMQCluster := &ConsumerRabbitMQ{
		Connection: conn,
		Channel:    channel,
	}

	nodeSequence = 0
	go func(cfg ConsumerConfig, consumerController controller.ConsumerController, seq *int) {
		for {
			reason, ok := <-rabbitMQCluster.Connection.NotifyClose(make(chan *amqp.Error))
			if !ok {
				log.Println("rabbitMQ cluster connection closed by termination signal")
				break
			}
			log.Println("rabbitMQ cluster connection closed, reason: ", reason)

			// reconnect with another node of cluster
			for {
				newSeq := next(cfg.BrokerClusterURLs, *seq)
				*seq = newSeq

				conn, err := amqp.Dial(cfg.BrokerClusterURLs[newSeq])
				if err == nil {
					rabbitMQCluster.mu.Lock()
					rabbitMQCluster.Connection = conn
					rabbitMQCluster.mu.Unlock()
					log.Println("rabbitMQ cluster reconnect success")
					// we can do re-open the channel directly without listening closing signal
					// because when connection is closed, the channel is automatically closed
					ch, err := conn.Channel()
					if err == nil {
						rabbitMQCluster.mu.Lock()
						rabbitMQCluster.Channel = ch
						rabbitMQCluster.mu.Unlock()
						// stop current consume process
						// rabbitMQCluster.StopConsuming()
						rabbitMQCluster.StartConsume(cfg, consumerController)
						log.Println("rabbitMQ cluster re-open channel success")
						break
					}

					log.Println("rabbitMQ cluster re-open channel failed, err: ", err)
					time.Sleep(time.Duration(1) * time.Second)
				}

				log.Println("rabbitMQ cluster reconnect failed, err: ", err)
				time.Sleep(time.Duration(1) * time.Second)
			}
		}
	}(cfg, consumerController, &nodeSequence)

	return rabbitMQCluster, nil
}

func (c *ConsumerRabbitMQ) StopConsuming() {
	close(c.stopConsuming)
}

func (c *ConsumerRabbitMQ) StartConsume(cfg ConsumerConfig, consumerController controller.ConsumerController) {
	c.stopConsuming = make(chan struct{})

	// start consume message
	go func(cluster *ConsumerRabbitMQ) {
		for {
			select {
			case <-c.stopConsuming:
				return
			case reason := <-c.Connection.NotifyClose(make(chan *amqp.Error)):
				log.Println("rabbitMQ cluster connection closed, reason: ", reason)
				return
			default:
				messages, err := cluster.Channel.Consume(cfg.QueueName, cfg.ConsumerName, false, false, false, false, nil)
				if err != nil {
					log.Fatalln("failed to consume message from producer, ", err)
				}

				for delivery := range messages {
					log.Println("[DEBUG] HALO")
					consumerController.Handle(context.Background(), delivery)
				}
			}
		}
	}(c)

	// start consume dropped message
	go func(cluster *ConsumerRabbitMQ) {
		for {
			select {
			case <-c.stopConsuming:
				return
			case reason := <-c.Connection.NotifyClose(make(chan *amqp.Error)):
				log.Println("rabbitMQ cluster connection closed, reason: ", reason)
				return
			default:
				droppedMessages, err := cluster.Channel.Consume(cfg.DLQName, cfg.ConsumerDroppedMessageName, false, false, false, false, nil)
				if err != nil {
					log.Fatalln("failed to consume dropped message, ", err)
				}

				for droppedDelivery := range droppedMessages {
					consumerController.HandleDroppedMessage(context.Background(), droppedDelivery)
				}
			}
		}
	}(c)
}

// Next element index of slice
func next(s []string, lastSeq int) int {
	length := len(s)
	if length == 0 || lastSeq == length-1 {
		return 0
	} else if lastSeq < length-1 {
		return lastSeq + 1
	} else {
	}
	return -1
}
