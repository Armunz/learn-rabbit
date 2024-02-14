package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/goombaio/namegenerator"
)

type Request struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	N := 1
	seed := time.Now().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	log.Println("===Start Injecting Data To Backend===")
	for i := 0; i < N; i++ {
		name := nameGenerator.Generate()
		age := rand.Intn(100)

		message := Request{
			Name: name,
			Age:  age,
		}

		if err := injectBackend(context.Background(), message); err != nil {
			log.Println("failed to inject request to backend, ", err)
			break
		}
		log.Println("===request send===")
		time.Sleep(5 * time.Second)
	}

	log.Println("===Finish===")
}

func injectBackend(ctx context.Context, message Request) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(5000)*time.Millisecond)
	defer cancel()

	requestBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctxTimeout, "POST", "http://localhost:8888/user", bytes.NewBuffer(requestBody))
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
