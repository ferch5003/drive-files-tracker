package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"user-service/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewConnection(config *config.EnvVars) (*sqlx.DB, error) {
	db := sqlx.MustConnect("pgx", config.PostgreDSN)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("DB Connected!")

	return db, nil
}
