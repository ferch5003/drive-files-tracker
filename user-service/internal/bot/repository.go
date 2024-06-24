package bot

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user-service/internal/domain"
)

// Queries.
const (
	_getBotByNameStmt = `SELECT id, name FROM bots WHERE name = ?;`
)

type Repository interface {
	// Get obtain one Bot by the name.
	Get(ctx context.Context, name string) (domain.Bot, error)
}

type repository struct {
	conn *sqlx.DB
}

func NewRepository(conn *sqlx.DB) Repository {
	return &repository{conn: conn}
}

func (r *repository) Get(ctx context.Context, name string) (domain.Bot, error) {
	var bot domain.Bot

	if err := r.conn.GetContext(ctx, &bot, _getBotByNameStmt, name); err != nil {
		return domain.Bot{}, err
	}

	return bot, nil
}
