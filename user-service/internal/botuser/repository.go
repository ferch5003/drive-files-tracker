package botuser

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"user-service/internal/domain"
)

// Queries.
const (
	_getAllParentsStmt = `SELECT bot_id, user_id, date, folder_id, is_parent
						  FROM bot_user
						  WHERE is_parent = TRUE;`
	_saveManyBotUserStmt = `INSERT INTO bot_user 
    						(bot_id, user_id, date, folder_id, is_parent) 
							VALUES %s;`
)

type Repository interface {
	// GetAllParents obtain all bot users from the database.
	GetAllParents(ctx context.Context) ([]domain.BotUser, error)

	// SaveMany generates a bulk insert to new entries for bot user.
	SaveMany(ctx context.Context, botUsers []domain.BotUser) error
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

func (r *repository) SaveMany(ctx context.Context, botUsers []domain.BotUser) error {
	tx, err := r.conn.Beginx()
	if err != nil {
		return err
	}

	var queryValues string
	var params []any
	for i, botUser := range botUsers {
		queryValues += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		params = append(params, botUser.BotID, botUser.UserID, botUser.Date, botUser.FolderID, botUser.IsParent)
	}

	queryValues = queryValues[:len(queryValues)-1]
	queryBulkInsert := fmt.Sprintf(_saveManyBotUserStmt, queryValues)

	stmt, err := tx.PreparexContext(ctx, queryBulkInsert)
	if err != nil {
		return err
	}

	defer func() {
		err = stmt.Close()
	}()

	res, err := stmt.ExecContext(ctx, params...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
