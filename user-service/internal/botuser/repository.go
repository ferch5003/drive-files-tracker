package botuser

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user-service/internal/domain"
)

// Queries.
const (
	_getAllParentsStmt = `SELECT bot_id, user_id, date, folder_id, is_parent
						  FROM bot_user
						  WHERE is_parent = TRUE;`
)

type Repository interface {
	// GetAllParents obtain all bot users from the database.
	GetAllParents(ctx context.Context) ([]domain.BotUser, error)
}

type repository struct {
	conn *sqlx.DB
}

func NewRepository(conn *sqlx.DB) Repository {
	return &repository{conn: conn}
}

func (r *repository) GetAllParents(ctx context.Context) ([]domain.BotUser, error) {
	botUser := make([]domain.BotUser, 0)

	if err := r.conn.SelectContext(ctx, &botUser, _getAllParentsStmt); err != nil {
		return make([]domain.BotUser, 0), err
	}

	return botUser, nil
}
