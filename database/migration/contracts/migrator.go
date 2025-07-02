package contracts

import (
	"context"

	migrationConfig "frisboo-bank/pkg/db/migration/config"
)

type Migrator interface {
	Up(ctx context.Context, version uint) error
	Down(ctx context.Context, version uint) error
	Config() *migrationConfig.MigrationOptions
}
