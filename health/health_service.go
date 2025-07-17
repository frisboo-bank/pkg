package health

import (
	"context"
	"sync"

	"frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/health/options"
)

type healthService struct {
	services   []contracts.HealthServiceCheck
	statusUp   contracts.StatusType
	statusDown contracts.StatusType
}

func (s *healthService) WithStatusDown(statusDown string) contracts.HealthService {
	s.statusDown = contracts.StatusType(s.statusDown)
	return s
}

func (s *healthService) WithStatusUp(statusUp string) contracts.HealthService {
	s.statusUp = contracts.StatusType(statusUp)
	return s
}

var _ contracts.HealthService = (*healthService)(nil)

func NewHealthService(
	services []contracts.HealthServiceCheck,
) contracts.HealthService {
	return &healthService{
		services:   services,
		statusUp:   options.StatusTypeUp,
		statusDown: options.StatusTypeDown,
	}
}

func (s *healthService) CheckHealth(ctx context.Context) contracts.CheckAllStatus {
	servicesCheck := make(contracts.CheckAllStatus)

	sync.OnceFunc(func() {
		for _, service := range s.services {
			servicesCheck[service.GetServiceName()] = contracts.NewCheckStatus(
				service.CheckHealth(ctx),
				s.statusUp,
				s.statusDown,
			)
		}
	})

	return servicesCheck
}
