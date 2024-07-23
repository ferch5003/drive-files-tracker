package botuser

import (
	"context"
	"user-service/internal/domain"
)

type Service interface {
	GetAllParents(ctx context.Context) ([]domain.BotUser, error)
}

type service struct {
	botUserRepository Repository
}

func NewService(botUserRepository Repository) Service {
	return &service{
		botUserRepository: botUserRepository,
	}
}

func (s service) GetAllParents(ctx context.Context) ([]domain.BotUser, error) {
	return s.botUserRepository.GetAllParents(ctx)
}
