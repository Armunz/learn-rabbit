package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"learn-rabbit/backend-service/internal/model"
	"net/http"
	"time"
)

type BackendRepository interface {
	SaveUser(ctx context.Context, request model.UserRequest) error
}

type BackendRepositoryImpl struct {
	producerURL  string
	exchangeName string
	routingKey   string
	timeoutMs    int
}

func NewBackendRepository(producerURL string, exchangeName string, routingKey string, timeoutMs int) BackendRepository {
	return &BackendRepositoryImpl{
		producerURL:  producerURL,
		exchangeName: exchangeName,
		routingKey:   routingKey,
		timeoutMs:    timeoutMs,
	}
}

// SaveUser implements BackendRepository
func (r *BackendRepositoryImpl) SaveUser(ctx context.Context, request model.UserRequest) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	message := Message{
		ExchangeName: r.exchangeName,
		RoutingKey:   r.routingKey,
		Data:         request,
	}

	requestBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctxTimeout, "POST", r.producerURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("expected 201 HTTP Status, but found %d", res.StatusCode)
	}

	return nil
}
