package user

import (
	"context"
	botdriven "user-service/internal/bot"
	"user-service/internal/domain"
)

type Service interface {
	// GetAll obtain all users.
	GetAll(ctx context.Context) ([]domain.User, error)

	// FindFolderID obtain the folder ID associated with a user and a bot.
	FindFolderID(ctx context.Context, userUsername, botName, date string) (string, error)

	// GetSpreadsheetData obtain the spreadsheet ID and GID associated with a user and a bot.
	GetSpreadsheetData(ctx context.Context, userUsername, botName, date string) (id, gid, column string, err error)
}

type service struct {
	userRepository Repository
	botRepository  botdriven.Repository
}

func NewService(userRepository Repository, botRepository botdriven.Repository) Service {
	return &service{
		userRepository: userRepository,
		botRepository:  botRepository,
	}
}

func (s service) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.userRepository.GetAll(ctx)
}

func (s service) FindFolderID(ctx context.Context, userUsername, botName, date string) (string, error) {
	user, err := s.userRepository.Get(ctx, userUsername)
	if err != nil {
		return "", err
	}

	bot, err := s.botRepository.Get(ctx, botName)
	if err != nil {
		return "", err
	}

	folderID, err := s.userRepository.FindFolderID(ctx, user.ID, bot.ID, date)
	if err != nil {
		return "", err
	}

	return folderID, nil
}

func (s service) GetSpreadsheetData(
	ctx context.Context,
	userUsername,
	botName,
	date string,
) (id, gid, column string, err error) {
	user, err := s.userRepository.Get(ctx, userUsername)
	if err != nil {
		return "", "", "", err
	}

	bot, err := s.botRepository.Get(ctx, botName)
	if err != nil {
		return "", "", "", err
	}

	id, gid, column, err = s.userRepository.GetSpreadsheetData(ctx, user.ID, bot.ID, date)
	if err != nil {
		return "", "", "", err
	}

	return id, gid, column, nil
}
