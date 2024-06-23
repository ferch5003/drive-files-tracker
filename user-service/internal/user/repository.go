package user

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user-service/internal/domain"
)

// Queries.
const (
	_getAllUsersStmt = `SELECT id, username FROM users;`
)

type Repository interface {
	// GetAll obtain all users from the database.
	GetAll(ctx context.Context) ([]domain.User, error)
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
