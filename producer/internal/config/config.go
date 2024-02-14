package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ProducerConfig struct {
	BrokerURL         string
	BrokerClusterURLs []string
	IsUsingCluster    bool
	ExchangeName      string
	ExchangeType      string
	APITimeoutMs      int
	RepoTimeoutMs     int
}

func InitConfig() ProducerConfig {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("failed to load config file, ", err)
	}

	brokerURL := os.Getenv("RABBITMQ_BROKER")
	exchangeName := os.Getenv("EXCHANGE_NAME")
	exchangeType := os.Getenv("EXCHANGE_TYPE")
	apiTimeout := os.Getenv("API_TIMEOUT_MS")
	repoTimeout := os.Getenv("REPO_TIMEOUT_MS")

	isUsingCluster, err := strconv.ParseBool(os.Getenv("IS_USING_CLUSTER"))
	if err != nil {
		log.Fatalln("failed to parse is using cluster, ", err)
	}

	clusterURLs := strings.Split(os.Getenv("RABBITMQ_CLUSTER"), "|")
	if len(clusterURLs) == 0 {
		log.Fatalln("cluster url must be specified")
	}

	apiTimeoutNum, err := strconv.Atoi(apiTimeout)
	if err != nil {
		log.Fatal("failed to parse api timeout, ", err)
	}

	repoTimeoutNum, err := strconv.Atoi(repoTimeout)
	if err != nil {
		log.Fatal("failed to parse repo timeout, ", err)
	}

	return ProducerConfig{
		BrokerURL:         brokerURL,
		BrokerClusterURLs: clusterURLs,
		IsUsingCluster:    isUsingCluster,
		ExchangeName:      exchangeName,
		ExchangeType:      exchangeType,
		APITimeoutMs:      apiTimeoutNum,
		RepoTimeoutMs:     repoTimeoutNum,
	}
}
