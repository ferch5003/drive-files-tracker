package botuser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func TestRepositorySaveMany_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	botUsers := []domain.BotUser{
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
	mock.ExpectBegin()
	mock.ExpectPrepare(`INSERT INTO bot_user`)
	mock.ExpectExec(`INSERT INTO bot_user`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository := NewRepository(dbx)

	// When
	err = repository.SaveMany(ctx, botUsers)

	// Then
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositorySaveMany_FailsDueToInvalidBeginTransaction(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	botUsers := []domain.BotUser{
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

	expectedError := errors.New("You have an error in your SQL syntax")

	mock.ExpectBegin().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.SaveMany(ctx, botUsers)

	// Then
	require.ErrorContains(t, err, "You have an error in your SQL syntax")
}

func TestRepositorySaveMany_FailsDueToInvalidPreparation(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	botUsers := []domain.BotUser{
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
	wrongQuery := regexp.QuoteMeta(`INSERT INTO bot_user (bot_id, user_id, date, folder_id, is_parent) 
										VALUES ($$, $$, $$, $$, $$);`)
	expectedError := errors.New(`Prepare: could not match actual sql: INSERT INTO bot_user`)

	mock.ExpectBegin()
	mock.ExpectPrepare(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.SaveMany(ctx, botUsers)

	// Then
	require.ErrorContains(t, err, "Prepare: could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}

func TestRepositorySaveMany_FailsDueToFailingExec(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	botUsers := []domain.BotUser{
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

	expectedError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")

	mock.ExpectBegin()
	mock.ExpectPrepare(`INSERT INTO bot_user`)
	mock.ExpectExec(`INSERT INTO bot_user`).WillReturnError(expectedError)
	mock.ExpectRollback()

	repository := NewRepository(dbx)

	// When
	err = repository.SaveMany(ctx, botUsers)

	// Then
	require.ErrorContains(t, err, "Error Code: 1136")
	require.ErrorContains(t, err, "Column count doesn't match value count at row 1")
}

func TestRepositorySaveMany_FailsDueToFailingExecWithFailingRollback(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	botUsers := []domain.BotUser{
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

	expectedExecError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")
	expectedRollbackError := fmt.Errorf("insert failed: %v, unable to back: %v",
		expectedExecError, "Rollack error")

	mock.ExpectBegin()
	mock.ExpectPrepare(`INSERT INTO bot_user`)
	mock.ExpectExec(`INSERT INTO bot_user`).WillReturnError(expectedExecError)
	mock.ExpectRollback().WillReturnError(expectedRollbackError)

	repository := NewRepository(dbx)

	// When
	err = repository.SaveMany(ctx, botUsers)

	// Then
	require.ErrorContains(t, err, "insert failed")
	require.ErrorContains(t, err, "Error Code: 1136")
	require.ErrorContains(t, err, "Column count doesn't match value count at row 1")
	require.ErrorContains(t, err, "unable to back")
	require.ErrorContains(t, err, "Rollack error")
}

func TestRepositorySaveMany_FailsDueToFailingCommit(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	botUsers := []domain.BotUser{
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
	expectedError := errors.New("sql: transaction has already been committed or rolled back")

	mock.ExpectBegin()
	mock.ExpectPrepare(`INSERT INTO bot_user`)
	mock.ExpectExec(`INSERT INTO bot_user`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.SaveMany(ctx, botUsers)

	// Then
	require.ErrorContains(t, err, "sql")
	require.ErrorContains(t, err, "transaction has already been committed or rolled back")
}
