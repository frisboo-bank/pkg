package postgresSqlx

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/db/postgres"
	"frisboo-bank/pkg/db/postgres/errors"

	"github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
)

type PostgresSqlx struct {
	DB               *sql.DB
	sqlxDB           *sqlx.DB
	statementBuilder squirrel.StatementBuilderType
	config           *postgres.PgConfig
}

func NewPostgresSqlx(cfg *postgres.PgConfig) (*PostgresSqlx, error) {
	maxOpenConns, _ := strconv.Atoi(os.Getenv(constants.DB_MAX_OPEN_CONNS_ENV))
	maxIdleConns, _ := strconv.Atoi(os.Getenv(constants.DB_MAX_IDLE_CONNS_ENV))
	connMaxLifetime, _ := strconv.Atoi(os.Getenv(constants.DB_CONN_MAX_FILETIME_ENV))

	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbName=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
	)

	db, err := sqlx.Connect(constants.DRIVER_NAME_POSTGRES, dataSourceName)
	if err != nil {
		panic(&errors.PostgresFailedToConnectError{Err: err})
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetime))

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("postgres: failed to ping the database and received the error: %v", err)
	}

	statementBuilder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(db)

	return &PostgresSqlx{
		DB:               db.DB,
		sqlxDB:           db,
		statementBuilder: statementBuilder,
		config:           cfg,
	}, nil
}

func (db *PostgresSqlx) Close() error {
	errDB := db.DB.Close()
	errSqlxDB := db.sqlxDB.Close()

	if errDB != nil {
		return fmt.Errorf("postgres: failed to close the DB connection with error: %v", errDB)
	}

	if errSqlxDB != nil {
		return fmt.Errorf("postgres: failed to close the Sqlx DB connection with error: %v", errSqlxDB)
	}

	return nil
}

func (db *PostgresSqlx) ExecuteTransaction(ctx context.Context, cb func(*PostgresSqlx) error) error {
	tx, err := db.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	err = cb(db)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf(
				"postgres: failed to execure transaction with error %v and the rollback also failed with error: %v",
				err,
				rollbackErr,
			)
		}

		return err
	}

	return tx.Commit()
}
