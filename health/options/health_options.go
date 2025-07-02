package options

import (
	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/reflection/typemapper"
	"net/http"

	"github.com/stoewer/go-strcase"
)

type HealthOptions struct {
	EndpointPath   string               `mapstructure:"endpointPath"`
	StatusCodeUp   uint                 `mapstructure:"statusCodeUp"`
	StatusUp       contracts.StatusType `mapstructure:"statusUp"`
	StatusCodeDown uint                 `mapstructure:"statusCodeDown"`
	StatusDown     contracts.StatusType `mapstructure:"statusDown"`
	Services       []contracts.HealthServiceCheck
}

type HealthOption = func(options *HealthOptions)

var DefaultHealthOptions = HealthOptions{
	EndpointPath:   "/healthz",
	StatusCodeUp:   http.StatusOK,
	StatusUp:       contracts.StatusTypeUp,
	StatusCodeDown: http.StatusServiceUnavailable,
	StatusDown:     contracts.StatusTypeDown,
}

func WithServices(services []contracts.HealthServiceCheck) HealthOption {
	return func(options *HealthOptions) {
		options.Services = services
	}
}

var optionName = strcase.LowerCamelCase(typemapper.GetGenericTypeNameByT[HealthOptions]())

func ProvideHealthOptions(env environment.Environment) (*HealthOptions, error) {
	opt, err := config.BindConfigKey[*HealthOptions](optionName, env)
	if err != nil {
		return nil, err
	}

	if opt.EndpointPath == "" {
		opt.EndpointPath = DefaultHealthOptions.EndpointPath
	}

	if opt.StatusCodeUp <= 0 {
		opt.StatusCodeUp = DefaultHealthOptions.StatusCodeUp
	}

	if opt.StatusUp == "" {
		opt.StatusUp = DefaultHealthOptions.StatusUp
	}

	if opt.StatusCodeDown <= 0 {
		opt.StatusCodeDown = DefaultHealthOptions.StatusCodeDown
	}

	if opt.StatusDown == "" {
		opt.StatusDown = DefaultHealthOptions.StatusDown
	}

	return opt, nil
}
