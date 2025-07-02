package postgresSqlx

import (
	"fmt"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/db/postgres"

	"github.com/jmoiron/sqlx"
)

type PostgresSqlx struct{}

func NewPostgresSqlx(cfg *postgres.PgConfig) (*PostgresSqlx, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbName=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sqlx.Connect(constants.DRIVER_NAME_POSTGRES, dataSourceName)
}

func createDB(cfg *postgres.PgConfig) error {
	if err != nil {
		return fmt.Errorf("postgres: failed to connect to the database with error: %v", err)
	}
}
