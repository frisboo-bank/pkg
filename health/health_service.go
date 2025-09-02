package health

import (
	"context"
	"sync"

	"frisboo-bank/pkg/health/config"
	"frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/syserrors"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

var _ contracts.HealthService = (*healthService)(nil)

type HealthServiceParams struct {
	Cfg      *config.Config
	Logger   loggerContracts.Logger
	Services []contracts.HealthServiceCheck
}

type healthService struct {
	cfg      *config.Config
	logger   loggerContracts.Logger
	services []contracts.HealthServiceCheck
}

func NewHealthService(params HealthServiceParams) contracts.HealthService {
	syserrors.AssertNotNil("Cfg", params.Cfg)
	syserrors.AssertNotNil("Logger", params.Logger)
	syserrors.AssertNotNil("Services", params.Services)

	return &healthService{
		cfg:      params.Cfg,
		logger:   params.Logger,
		services: params.Services,
	}
}

func (s *healthService) AddServices(services ...contracts.HealthServiceCheck) {
	s.services = services
}

func (s *healthService) CheckHealth(ctx context.Context) contracts.CheckAllStatus {
	servicesCheck := make(contracts.CheckAllStatus)

	sync.OnceFunc(func() {
		for _, service := range s.services {
			servicesCheck[service.GetServiceName()] = contracts.NewCheckStatus(
				service.CheckHealth(ctx),
				s.cfg.StatusUp,
				s.cfg.StatusDown,
			)
		}
	})

	return servicesCheck
}

func (s *healthService) Logger() loggerContracts.Logger {
	return s.logger
}
