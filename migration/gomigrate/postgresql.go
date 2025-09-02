package gomigrate

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type goMigratePostgresqlMigrator struct {
	// config     *migrationConfig.MigrationOptions
	// datasource string
	// migrate    *migrate.Migrate
}

// func NewPostgresqlMigrator(config *migrationConfig.MigrationOptions) (contracts.Migrator, error) {
// 	if strings.TrimSpace(config.DBName) == "" {
// 		return nil, syserrors.Newf("DBName is required in the config")
// 	}
//
// 	if strings.TrimSpace(config.MigrationDir) == "" {
// 		return nil, syserrors.Newf("MigrationDir is required in the config")
// 	}
//
// 	datasource := postgres.GenerateDataSource(&postgres.PgConfig{
// 		DBName:   config.DBName,
// 		Host:     config.Host,
// 		Password: config.Password,
// 		Port:     config.Port,
// 		SSLMode:  config.SSLMode,
// 		User:     config.User,
// 	})
//
// 	migrate, err := migrate.New(fmt.Sprintf("file://%s", config.MigrationDir), datasource)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &goMigratePostgresqlMigrator{
// 		config:     config,
// 		datasource: datasource,
// 		migrate:    migrate,
// 	}, nil
// }
//
// func (g *goMigratePostgresqlMigrator) Down(ctx context.Context, version uint) error {
// 	return g.executeCommand(migrationConfig.CommandTypeDown, version)
// }
//
// func (g *goMigratePostgresqlMigrator) Up(ctx context.Context, version uint) error {
// 	return g.executeCommand(migrationConfig.CommandTypeUp, version)
// }
//
// func (g *goMigratePostgresqlMigrator) Config() *migrationConfig.MigrationOptions {
// 	return g.config
// }
//
// func (g *goMigratePostgresqlMigrator) executeCommand(commandType migrationConfig.CommandType, version uint) error {
// 	var err error
//
// 	switch true {
// 	case version > 0:
// 		err = g.migrate.Migrate(version)
// 	case commandType == migrationConfig.CommandTypeUp:
// 		err = g.migrate.Up()
// 	case commandType == migrationConfig.CommandTypeDown:
// 		err = g.migrate.Down()
// 	}
//
// 	if errors.Is(err, migrate.ErrNoChange) {
// 		return nil
// 	}
//
// 	return err
// }
