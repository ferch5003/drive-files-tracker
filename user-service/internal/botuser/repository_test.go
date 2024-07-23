package botuser

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

func TestRepositoryGetAllParents_Successful(t *testing.T) {
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

	expectedBotUsers := []domain.BotUser{
		{
			BotID:    1,
			UserID:   1,
			Date:     "test",
			FolderID: "test",
			IsParent: false,
		},
		{
			BotID:    1,
			UserID:   2,
			Date:     "test",
			FolderID: "test",
			IsParent: true,
		},
	}

	columns := []string{"bot_id", "user_id", "date", "folder_id", "is_parent"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(
		expectedBotUsers[0].BotID,
		expectedBotUsers[0].UserID,
		expectedBotUsers[0].Date,
		expectedBotUsers[0].FolderID,
		expectedBotUsers[0].IsParent)
	rows.AddRow(
		expectedBotUsers[1].BotID,
		expectedBotUsers[1].UserID,
		expectedBotUsers[1].Date,
		expectedBotUsers[1].FolderID,
		expectedBotUsers[1].IsParent)
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	repository := NewRepository(dbx)

	// When
	botUsers, err := repository.GetAllParents(ctx)

	// Then
	require.NoError(t, err)
	require.NotNil(t, botUsers)
	require.Equal(t, expectedBotUsers, botUsers)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetAllParents_FailsDueToInvalidSelect(t *testing.T) {
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

	wrongQuery := regexp.QuoteMeta("SELECT wrong FROM bot_user;")
	expectedError := errors.New(`Query: could not match actual sql: \"SELECT\"`)

	expectedBotUsers := make([]domain.BotUser, 0)
	mock.ExpectQuery(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	users, err := repository.GetAllParents(ctx)

	// Then
	require.Equal(t, expectedBotUsers, users)
	require.ErrorContains(t, err, "Query")
	require.ErrorContains(t, err, "could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}
