package mongodb

import (
	"context"
	"database/sql"
	"fmt"

	"frisboo-bank/pkg/database/database_client/config"
	"frisboo-bank/pkg/database/database_client/contracts"
	databaseclienttype "frisboo-bank/pkg/database/database_client/enums/database_client_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/validation"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var _ contracts.DatabaseClientAdapter = (*mongoDBDatabaseClientAdapter)(nil)

type mongoDBDatabaseClientAdapter struct {
	name   string
	cfg    *config.Config
	client *mongo.Client
	ctx    context.Context
	logger loggerContracts.Logger
	// meter        metric.Meter
}

func New(name string, cfg *config.Config, logger loggerContracts.Logger) (contracts.DatabaseClientAdapter, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	mongoOpts := mongoOptions.Client().
		ApplyURI(uri).
		SetAuth(mongoOptions.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		}).
		SetConnectTimeout(cfg.ConnectionTimeout).
		SetMaxConnIdleTime(cfg.MaxConnectionIdleTime).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxPoolSize(cfg.MaxPoolSize)

	ctx := context.Background()
	if cfg.Context != nil {
		ctx = cfg.Context
	}

	if cfg.EnableTracing {
		mongoOpts.Monitor = otelmongo.NewMonitor()
	}

	client, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		return nil, err
	}

	if err := mgm.SetDefaultConfig(nil, cfg.Database, mongoOpts); err != nil {
		return nil, err
	}

	return &mongoDBDatabaseClientAdapter{
		name:   name,
		cfg:    cfg,
		client: client,
		ctx:    ctx,
		logger: logger,
	}, nil
}

func (m *mongoDBDatabaseClientAdapter) Ping() error {
	return m.client.Ping(m.ctx, nil)
}

func (m *mongoDBDatabaseClientAdapter) Disconnect() error {
	return m.client.Disconnect(m.ctx)
}

func (m *mongoDBDatabaseClientAdapter) Name() string {
	return m.name
}

func (m *mongoDBDatabaseClientAdapter) Type() databaseclienttype.DatabaseClientType {
	return databaseclienttype.DatabaseClientTypes.MONGODB
}

func (m *mongoDBDatabaseClientAdapter) Config() *config.Config {
	return m.cfg
}

func (p *mongoDBDatabaseClientAdapter) DB() *sql.DB {
	return nil
}

func (m *mongoDBDatabaseClientAdapter) Logger() loggerContracts.Logger {
	return m.logger
}
