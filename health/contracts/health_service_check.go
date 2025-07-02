package contracts

import "context"

type HealthServiceCheck interface {
	CheckHealth(ctx context.Context) error
	GetServiceName() string
}
