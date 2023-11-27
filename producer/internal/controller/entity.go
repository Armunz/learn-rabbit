package controller

type Request struct {
	ExchangeName string      `json:"exchange_name"`
	RoutingKey   string      `json:"routing_key"`
	Data         interface{} `json:"data"`
}
