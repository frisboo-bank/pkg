package contracts

import (
	"context"
)

type Migrator interface {
	Up(ctx context.Context, version uint) error
	Down(ctx context.Context, version uint) error
}
