package user

import (
	"context"
	"user-service/internal/domain"
)

type Service interface {
	// GetAll obtain all users.
	GetAll(ctx context.Context) ([]domain.User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s service) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.repository.GetAll(ctx)
}
