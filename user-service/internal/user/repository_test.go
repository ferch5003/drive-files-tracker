package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"user-service/internal/domain"
)

func TestRepositoryGetAll_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			require.Error(t, err)
		}
	}(db)

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUsers := []domain.User{
		{
			ID:       1,
			Username: "john",
		},
		{
			ID:       2,
			Username: "jane",
		},
	}

	columns := []string{"id", "username"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(expectedUsers[0].ID, expectedUsers[0].Username)
	rows.AddRow(expectedUsers[1].ID, expectedUsers[1].Username)
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	repository := NewRepository(dbx)

	// When
	users, err := repository.GetAll(ctx)

	// Then
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Equal(t, expectedUsers, users)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetAll_FailsDueToInvalidSelect(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			require.Error(t, err)
		}
	}(db)

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	wrongQuery := regexp.QuoteMeta("SELECT wrong FROM users;")
	expectedError := errors.New(`Query: could not match actual sql: \"SELECT first_name, last_name,
										email FROM users;\" with expected regexp \"SELECT wrong FROM users;\"`)

	expectedUsers := make([]domain.User, 0)
	mock.ExpectQuery(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	users, err := repository.GetAll(ctx)

	// Then
	require.Equal(t, expectedUsers, users)
	require.ErrorContains(t, err, "Query")
	require.ErrorContains(t, err, "could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}

func TestRepositoryFindFolderID_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedFolderID := "Test"

	user := domain.User{
		ID:       1,
		Username: "Test User",
	}

	bot := domain.Bot{
		ID:   1,
		Name: "Test Bot",
	}

	columns := []string{"folder_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(expectedFolderID)
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	repository := NewRepository(dbx)

	// When
	folderID, err := repository.FindFolderID(ctx, user.ID, bot.ID, "")

	// Then
	require.NoError(t, err)
	require.NotNil(t, bot)
	require.Equal(t, expectedFolderID, folderID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryFindFolderID_FailsDueToInvalidGet(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	wrongQuery := regexp.QuoteMeta("SELECT wrong FROM bot_user;")
	expectedError := errors.New(`Query: could not match actual sql: \"SELECT folder_id
									    FROM bot_user
									    INNER JOIN bots
										ON bot_user.bot_id = ?
										INNER JOIN users
 										ON bot_user.user_id = ?
						 				WHERE bot_user.date = ?;\" with expected regexp \"SELECT 
										wrong FROM bot_user;\"`)

	expectedFolderID := ""

	mock.ExpectQuery(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	folderID, err := repository.FindFolderID(ctx, 0, 0, "")

	// Then
	require.Equal(t, expectedFolderID, folderID)
	require.ErrorContains(t, err, "Query")
	require.ErrorContains(t, err, "could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}
