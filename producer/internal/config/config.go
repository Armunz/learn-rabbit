package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ProducerConfig struct {
	BrokerURL     string
	QueueName     string
	ExchangeName  string
	ExchangeType  string
	APITimeoutMs  int
	RepoTimeoutMs int
}

func InitConfig() ProducerConfig {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("failed to load config file, ", err)
	}

	brokerURL := os.Getenv("RABBITMQ_BROKER")
	queueName := os.Getenv("QUEUE_NAME")
	exchangeName := os.Getenv("EXCHANGE_NAME")
	exchangeType := os.Getenv("EXCHANGE_TYPE")
	apiTimeout := os.Getenv("API_TIMEOUT_MS")
	repoTimeout := os.Getenv("REPO_TIMEOUT_MS")

	apiTimeoutNum, err := strconv.Atoi(apiTimeout)
	if err != nil {
		log.Fatal("failed to parse api timeout, ", err)
	}

	repoTimeoutNum, err := strconv.Atoi(repoTimeout)
	if err != nil {
		log.Fatal("failed to parse repo timeout, ", err)
	}

	return ProducerConfig{
		BrokerURL:     brokerURL,
		QueueName:     queueName,
		ExchangeName:  exchangeName,
		ExchangeType:  exchangeType,
		APITimeoutMs:  apiTimeoutNum,
		RepoTimeoutMs: repoTimeoutNum,
	}
}
