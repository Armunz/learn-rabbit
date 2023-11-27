package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ProducerURL string
	TopicName   string
	RoutingKey  string
	APITimeout  int
	RepoTimeout int
}

func InitConfig() Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("failed to load config file, ", err)
	}

	producerURL := os.Getenv("PRODUCER_URL")
	exchangeName := os.Getenv("EXCHANGE_NAME")
	routingKey := os.Getenv("ROUTING_KEY")
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

	return Config{
		ProducerURL: producerURL,
		TopicName:   exchangeName,
		RoutingKey:  routingKey,
		APITimeout:  apiTimeoutNum,
		RepoTimeout: repoTimeoutNum,
	}
}
