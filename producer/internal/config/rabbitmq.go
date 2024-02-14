package config

import (
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ProducerRabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	mu         sync.Mutex
}

func InitRabbitMQ(cfg ProducerConfig) *ProducerRabbitMQ {
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

	return &ProducerRabbitMQ{
		Connection: connection,
		Channel:    channel,
	}
}

func InitRabbitMQCluster(cfg ProducerConfig) *ProducerRabbitMQ {
	rabbitMQCluster, err := DialCluster(cfg.BrokerClusterURLs)
	if err != nil {
		log.Fatalln("failed to connect to rabbitMQ cluster, ", err)
	}
	// rabbitMQCluster.OpenChannel()

	// declare exchange
	err = rabbitMQCluster.Channel.ExchangeDeclare(cfg.ExchangeName, cfg.ExchangeType, true, false, false, false, nil)
	if err != nil {
		log.Fatalln("failed to declare exchange, ", err)
	}

	return rabbitMQCluster
}

// based on https://github.com/sirius1024/go-amqp-reconnect/tree/master
func DialCluster(urls []string) (*ProducerRabbitMQ, error) {
	nodeSequence := 0
	var conn *amqp.Connection
	var err error
	for {
		conn, err = amqp.Dial(urls[nodeSequence])
		if err == nil {
			break
		}
		newSeq := next(urls, nodeSequence)
		nodeSequence = newSeq
		time.Sleep(time.Duration(1) * time.Second)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	rabbitMQCluster := &ProducerRabbitMQ{
		Connection: conn,
		Channel:    channel,
	}

	nodeSequence = 0
	go func(urls []string, seq *int) {
		for {
			reason, ok := <-rabbitMQCluster.Connection.NotifyClose(make(chan *amqp.Error))
			if !ok {
				log.Println("rabbitMQ cluster connection closed by termination signal")
				break
			}
			log.Println("rabbitMQ cluster connection closed, reason: ", reason)

			// reconnect with another node of cluster
			for {
				newSeq := next(urls, *seq)
				*seq = newSeq

				conn, err := amqp.Dial(urls[newSeq])
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
	}(urls, &nodeSequence)

	return rabbitMQCluster, nil
}

// func (p *ProducerRabbitMQ) OpenChannel() {
// 	channel, err := p.Connection.Channel()
// 	if err != nil {
// 		log.Fatalln("failed to open rabbitMQ cluster channel, ", err)
// 	}

// 	p.Channel = channel
// 	go func() {
// 		for {
// 			reason, ok := <-p.Channel.NotifyClose(make(chan *amqp.Error))
// 			// exit this goroutine if closed by developer
// 			if !ok || p.Channel.IsClosed() {
// 				log.Println("rabbitMQ cluster channel closed, reason: ", reason)
// 				p.Channel.Close() // close again, ensure closed flag set when connection closed
// 			}

// 			// reconnect if not closed by developer
// 			for {
// 				// wait 1s for connection reconnect
// 				time.Sleep(time.Duration(1) * time.Second)
// 				ch, err := p.Connection.Channel()
// 				if err == nil {
// 					p.mu.Lock()
// 					p.Channel = ch
// 					p.mu.Unlock()
// 					log.Println("rabbitMQ cluster re-open channel success")
// 					break
// 				}

// 				log.Println("rabbitMQ cluster re-open channel failed, err: ", err)
// 			}
// 		}
// 	}()
// }

// Next element index of slice
func next(s []string, lastSeq int) int {
	length := len(s)
	if length == 0 || lastSeq == length-1 {
		return 0
	} else if lastSeq < length-1 {
		return lastSeq + 1
	} else {
		return -1
	}
}
