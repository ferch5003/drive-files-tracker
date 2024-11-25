package user

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user-service/internal/domain"
)

// Queries.
const (
	_getAllUsersStmt       = `SELECT id, username FROM users;`
	_getUserByUsernameStmt = `SELECT id, username FROM users WHERE username = $1;`
	_findFolderIDStmt      = `SELECT DISTINCT folder_id
							  FROM bot_user
							  INNER JOIN bots
							  ON bot_user.bot_id = $1
							  INNER JOIN users
							  ON bot_user.user_id = $2
							  WHERE bot_user.date = $3;`
	_getSpreadsheetDataStmt = `SELECT 
    						   bot_user.spreadsheet_id, bot_user.spreadsheet_gid, bot_user.spreadsheet_column
							   FROM bot_user
							   INNER JOIN bots
							   ON bot_user.bot_id = bots.id
						       INNER JOIN users
							   ON bot_user.user_id = users.id
						       WHERE bot_user.bot_id = $1
							   AND bot_user.user_id = $2
						       AND bot_user.date = $3;`
)

type Repository interface {
	// GetAll obtain all users from the database.
	GetAll(ctx context.Context) ([]domain.User, error)

	// Get obtain one User by the username.
	Get(ctx context.Context, username string) (domain.User, error)

	// FindFolderID obtain the folder ID associated with a user and a bot.
	FindFolderID(ctx context.Context, userID, botID int, date string) (string, error)

	// GetSpreadsheetData obtain the spreadsheet ID and GID associated with a user and a bot.
	GetSpreadsheetData(ctx context.Context, userID, botID int, date string) (id, gid, column string, err error)
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

func (r *repository) Get(ctx context.Context, username string) (domain.User, error) {
	var user domain.User

	if err := r.conn.GetContext(ctx, &user, _getUserByUsernameStmt, username); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *repository) FindFolderID(ctx context.Context, userID, botID int, date string) (string, error) {
	var folderID string

	if err := r.conn.GetContext(ctx, &folderID, _findFolderIDStmt, botID, userID, date); err != nil {
		return "", err
	}

	return folderID, nil
}

func (r *repository) GetSpreadsheetData(
	ctx context.Context,
	userID,
	botID int,
	date string,
) (id, gid, column string, err error) {
	var data struct {
		SpreadsheetID     string `db:"spreadsheet_id"`
		SpreadsheetGID    string `db:"spreadsheet_gid"`
		SpreadsheetColumn string `db:"spreadsheet_column"`
	}

	if err := r.conn.GetContext(ctx, &data, _getSpreadsheetDataStmt, botID, userID, date); err != nil {
		return "", "", "", err
	}

	return data.SpreadsheetID, data.SpreadsheetGID, data.SpreadsheetColumn, nil
}
