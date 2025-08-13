package contracts

import (
	"context"
)

type HealthService interface {
	AddServices(services ...HealthServiceCheck)
	CheckHealth(ctx context.Context) CheckAllStatus
}
