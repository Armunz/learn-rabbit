package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConsumerConfig struct {
	BrokerURL       string
	DBUserDSN       string
	QueueName       string
	DLQName         string
	QueueRoutingKey string
	ExchangeName    string
	ExchangeType    string
	DLXName         string
	DLXType         string
	ConsumerName    string

	MYSQLDSN            string
	MYSQLMaxConn        int
	MYSQLIdleConn       int
	MYSQLQueryTimeoutMs int
	ConnLifeTimeSecond  int
	APITimeoutMs        int
	RepoTimeoutMs       int
}

func InitConfig() ConsumerConfig {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("failed to load config file, ", err)
	}

	// rabbitmq
	brokerURL := os.Getenv("RABBITMQ_BROKER")
	queueName := os.Getenv("QUEUE_NAME")
	dlqName := os.Getenv("DLQ_NAME")
	queueRoutingKey := os.Getenv("ROUTING_KEY")
	exchangeName := os.Getenv("EXCHANGE_NAME")
	exchangeType := os.Getenv("EXCHANGE_TYPE")
	dlxName := os.Getenv("DLX_NAME")
	dlxType := os.Getenv("DLX_TYPE")
	consumerName := os.Getenv("CONSUMER_NAME")

	// mysql
	mysqlDSN := fmt.Sprintf("%s:%s@%s?%s", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_URL"), os.Getenv("MYSQL_CONN_PARAMS"))
	maxConn := os.Getenv("MYSQL_CONN_MAX")
	idleConn := os.Getenv("MYSQL_CONN_IDLE")
	connLifeTimeSecond := os.Getenv("MYSQL_CONN_TIMEOUT_SECOND")
	queryTimeout := os.Getenv("MYSQL_QUERY_TIMEOUT_MS")

	maxConnNum, err := strconv.Atoi(maxConn)
	if err != nil {
		log.Fatalln("failed to parse mysql max connection, ", err)
	}

	idleConnNum, err := strconv.Atoi(idleConn)
	if err != nil {
		log.Fatalln("failed to parse mysql idle connection, ", err)
	}

	connLifeTimeSecondNum, err := strconv.Atoi(connLifeTimeSecond)
	if err != nil {
		log.Fatalln("failed to parse mysql connection lifetime, ", err)
	}

	queryTimeoutNum, err := strconv.Atoi(queryTimeout)
	if err != nil {
		log.Fatalln("failed to parse mysql query timeout, ", err)
	}

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

	return ConsumerConfig{
		BrokerURL:           brokerURL,
		DBUserDSN:           mysqlDSN,
		QueueName:           queueName,
		DLQName:             dlqName,
		QueueRoutingKey:     queueRoutingKey,
		ExchangeName:        exchangeName,
		ExchangeType:        exchangeType,
		DLXName:             dlxName,
		DLXType:             dlxType,
		ConsumerName:        consumerName,
		MYSQLDSN:            mysqlDSN,
		MYSQLMaxConn:        maxConnNum,
		MYSQLIdleConn:       idleConnNum,
		MYSQLQueryTimeoutMs: queryTimeoutNum,
		ConnLifeTimeSecond:  connLifeTimeSecondNum,
		APITimeoutMs:        apiTimeoutNum,
		RepoTimeoutMs:       repoTimeoutNum,
	}
}
