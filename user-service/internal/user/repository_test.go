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
