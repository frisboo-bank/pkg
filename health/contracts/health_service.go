package contracts

import (
	"context"
)

type HealthService interface {
	WithStatusUp(statusUp string) HealthService
	WithStatusDown(statusDown string) HealthService

	CheckHealth(ctx context.Context) CheckAllStatus
}
