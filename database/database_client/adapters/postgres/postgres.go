package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"frisboo-bank/pkg/database/database_client/config"
	"frisboo-bank/pkg/database/database_client/contracts"
	databaseclienttype "frisboo-bank/pkg/database/database_client/enums/database_client_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/validation"

	"github.com/jmoiron/sqlx"
)

var _ contracts.DatabaseClientAdapter = (*postgresDatabaseClientAdapter)(nil)

type postgresDatabaseClientAdapter struct {
	name   string
	client *sqlx.DB
	db     *sql.DB
	cfg    *config.Config
	ctx    context.Context
	logger loggerContracts.Logger
	// meter        metric.Meter
}

func New(name string, cfg *config.Config, logger loggerContracts.Logger) (contracts.DatabaseClientAdapter, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	sslOption := "enable"
	if !cfg.SSLMode {
		sslOption = "disable"
	}

	uri := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		sslOption,
	)

	ctx := context.Background()
	if cfg.Context != nil {
		ctx = cfg.Context
	}

	client, err := sqlx.ConnectContext(ctx, "pgx", uri)
	if err != nil {
		return nil, err
	}

	client.SetMaxIdleConns(int(cfg.MinPoolSize))
	client.SetMaxOpenConns(int(cfg.MaxPoolSize))
	client.SetConnMaxIdleTime(cfg.MaxConnectionIdleTime)

	return &postgresDatabaseClientAdapter{
		name:   name,
		client: client,
		db:     client.DB,
		cfg:    cfg,
		ctx:    ctx,
		logger: logger,
	}, nil
}

func (p *postgresDatabaseClientAdapter) Config() *config.Config {
	return p.cfg
}

func (p *postgresDatabaseClientAdapter) Disconnect() error {
	if err := p.client.Close(); err != nil {
		return err
	}
	return nil
}

func (p *postgresDatabaseClientAdapter) Logger() loggerContracts.Logger {
	return p.logger
}

func (p *postgresDatabaseClientAdapter) Name() string {
	return p.name
}

func (p *postgresDatabaseClientAdapter) Ping() error {
	return p.client.PingContext(p.ctx)
}

func (p *postgresDatabaseClientAdapter) Type() databaseclienttype.DatabaseClientType {
	return databaseclienttype.DatabaseClientTypes.POSTGRES
}
