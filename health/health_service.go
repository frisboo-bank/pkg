package health

import (
	"context"
	"sync"

	"frisboo-bank/pkg/customerrors"
	"frisboo-bank/pkg/health/config"
	"frisboo-bank/pkg/health/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/utils"
)

var _ contracts.HealthService = (*healthService)(nil)

var pHSError = customerrors.PrefixedError("health service")

type healthService struct {
	cfg      *config.Config
	logger   loggerContracts.Logger
	services []contracts.HealthServiceCheck
}

func NewHealthService(
	logger loggerContracts.Logger,
	opts *options.OptionBuilder[config.Config],
) contracts.HealthService {
	utils.Assert(logger != nil, pHSError.New("logger can't be nil"))
	utils.Assert(opts != nil, pHSError.New("opts can't be nil"))

	cfg := opts.Build()

	return &healthService{
		cfg:    cfg,
		logger: logger,
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
