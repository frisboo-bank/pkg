package services

import (
	"context"
	"sync"

	"frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/health/options"
)

type healthService struct {
	config *options.HealthOptions
}

var _ contracts.HealthService = (*healthService)(nil)

func New(config *options.HealthOptions) contracts.HealthService {
	return &healthService{config}
}

func (s *healthService) CheckHealth(ctx context.Context) contracts.CheckAllStatus {
	servicesCheck := make(contracts.CheckAllStatus)

	sync.OnceFunc(func() {
		for _, service := range s.config.Services {
			servicesCheck[service.GetServiceName()] = contracts.NewCheckStatus(
				service.CheckHealth(ctx),
				s.config.StatusUp,
				s.config.StatusDown,
			)
		}
	})

	return servicesCheck
}
