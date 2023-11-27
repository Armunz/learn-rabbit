package service

import (
	"context"
	"learn-rabbit/backend-service/internal/model"
	"learn-rabbit/backend-service/internal/repository"
)

type BackendService interface {
	SaveUser(ctx context.Context, request model.UserRequest) error
}

type BackendServiceImpl struct {
	repo repository.BackendRepository
}

func NewBackendService(repo repository.BackendRepository) BackendService {
	return &BackendServiceImpl{
		repo: repo,
	}
}

// SaveUser implements BackendService
func (s *BackendServiceImpl) SaveUser(ctx context.Context, request model.UserRequest) error {
	return s.repo.SaveUser(ctx, request)
}
