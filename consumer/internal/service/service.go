package service

import (
	"context"
	"encoding/json"
	"learn-rabbit/consumer/internal/repository"
	"log"
)

type consumerServiceImpl struct {
	userRepo repository.UserRepository
}

type ConsumerService interface {
	Handle(ctx context.Context, body []byte) error
}

func NewConsumerService(userRepo repository.UserRepository) ConsumerService {
	return &consumerServiceImpl{
		userRepo: userRepo,
	}
}

// Handle implements ConsumerService.
func (s *consumerServiceImpl) Handle(ctx context.Context, body []byte) error {
	var user repository.User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Println("failed to parse body to user, ", err)
		return err
	}

	if err := s.userRepo.SaveUser(ctx, user); err != nil {
		log.Println("failed to save user, ", err)
		return err
	}

	return nil
}
