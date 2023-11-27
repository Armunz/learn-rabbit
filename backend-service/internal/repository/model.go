package repository

type Message struct {
	ExchangeName string      `json:"exchange_name"`
	RoutingKey   string      `json:"routing_key"`
	Data         interface{} `json:"data"`
}
