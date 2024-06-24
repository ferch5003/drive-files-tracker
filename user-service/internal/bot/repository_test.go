package bot

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"user-service/internal/domain"
)

func TestRepositoryGet_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedBotName := "Test"
	expectedBot := domain.Bot{
		ID:   1,
		Name: "Test",
	}

	columns := []string{"id", "name"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(expectedBot.ID, expectedBot.Name)
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	repository := NewRepository(dbx)

	// When
	bot, err := repository.Get(ctx, expectedBotName)

	// Then
	require.NoError(t, err)
	require.NotNil(t, bot)
	require.Equal(t, expectedBot, bot)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGet_FailsDueToInvalidGet(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	wrongQuery := regexp.QuoteMeta("SELECT wrong FROM bots;")
	expectedError := errors.New(`Query: could not match actual sql: \"SELECT id, name 
										FROM bots WHERE name = ?;\" with expected regexp \"SELECT 
										wrong FROM bots;\"`)

	expectedBot := domain.Bot{}

	mock.ExpectQuery(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	bot, err := repository.Get(ctx, "")

	// Then
	require.Equal(t, expectedBot, bot)
	require.ErrorContains(t, err, "Query")
	require.ErrorContains(t, err, "could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}
