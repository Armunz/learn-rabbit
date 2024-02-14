package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"learn-rabbit/backend-service/internal/model"
	"log"
	"net/http"
	"time"
)

type BackendRepository interface {
	SaveUser(ctx context.Context, request model.UserRequest) error
}

type backendRepositoryImpl struct {
	producerURL  string
	exchangeName string
	routingKey   string
	timeoutMs    int
}

func NewBackendRepository(producerURL string, exchangeName string, routingKey string, timeoutMs int) BackendRepository {
	return &backendRepositoryImpl{
		producerURL:  producerURL,
		exchangeName: exchangeName,
		routingKey:   routingKey,
		timeoutMs:    timeoutMs,
	}
}

// SaveUser implements BackendRepository
func (r *backendRepositoryImpl) SaveUser(ctx context.Context, request model.UserRequest) error {
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

	req, err := http.NewRequestWithContext(ctxTimeout, "POST", r.producerURL+"/user", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("failed to create request, ", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("failed to do request, ", err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("failed to read response body, ", err)
	}

	log.Println("Response Body: ", string(body))

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("expected 201 HTTP Status, but found %d", res.StatusCode)
	}

	return nil
}
