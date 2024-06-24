package user

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user-service/internal/domain"
)

// Queries.
const (
	_getAllUsersStmt  = `SELECT id, username FROM users;`
	_findFolderIDStmt = `SELECT folder_id
						 FROM bot_user
						 INNER JOIN bots
						 ON bot_user.bot_id = ?
						 INNER JOIN users
						 ON bot_user.user_id = ?
						 WHERE bot_user.date = ?;`
)

type Repository interface {
	// GetAll obtain all users from the database.
	GetAll(ctx context.Context) ([]domain.User, error)

	// FindFolderID obtain the folder ID associated with a user and a bot.
	FindFolderID(ctx context.Context, userID, botID int, date string) (string, error)
}

type repository struct {
	conn *sqlx.DB
}

func NewRepository(conn *sqlx.DB) Repository {
	return &repository{conn: conn}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.User, error) {
	users := make([]domain.User, 0)

	if err := r.conn.SelectContext(ctx, &users, _getAllUsersStmt); err != nil {
		return make([]domain.User, 0), err
	}

	return users, nil
}

func (r *repository) FindFolderID(ctx context.Context, userID, botID int, date string) (string, error) {
	var folderID string

	if err := r.conn.GetContext(ctx, &folderID, _findFolderIDStmt, userID, botID, date); err != nil {
		return "", err
	}

	return folderID, nil
}
