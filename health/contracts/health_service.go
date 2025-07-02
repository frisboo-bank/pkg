package contracts

import (
	"context"
)

type HealthService interface {
	CheckHealth(ctx context.Context) CheckAllStatus
}
