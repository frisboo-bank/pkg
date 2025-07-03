package options

import (
	"fmt"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/options"
	optionsContracts "frisboo-bank/pkg/options/contracts"
	"net/http"
)

type HealthOptions struct {
	EndpointPath   string `mapstructure:"endpointPath"`
	StatusCodeUp   uint   `mapstructure:"statusCodeUp"`
	StatusUp       string `mapstructure:"statusUp"`
	StatusCodeDown uint   `mapstructure:"statusCodeDown"`
	StatusDown     string `mapstructure:"statusDown"`
	Services       []contracts.HealthServiceCheck
}

var defaultHealthOptions = HealthOptions{
	EndpointPath:   "/healthz",
	StatusCodeUp:   http.StatusOK,
	StatusUp:       contracts.StatusTypeUp,
	StatusCodeDown: http.StatusServiceUnavailable,
	StatusDown:     contracts.StatusTypeDown,
}

var _ optionsContracts.Options = (*HealthOptions)(nil)

func (o *HealthOptions) Clone() optionsContracts.Options {
	return options.CloneOptions(o)
}

func (o *HealthOptions) SetDefaults() {
	if o.EndpointPath == "" {
		o.EndpointPath = defaultHealthOptions.EndpointPath
	}

	if o.StatusCodeUp == 0 {
		o.StatusCodeUp = defaultHealthOptions.StatusCodeUp
	}

	if o.StatusUp == "" {
		o.StatusUp = defaultHealthOptions.StatusUp
	}

	if o.StatusCodeDown == 0 {
		o.StatusCodeDown = defaultHealthOptions.StatusCodeDown
	}

	if o.StatusDown == "" {
		o.StatusDown = defaultHealthOptions.StatusDown
	}
}

func (o *HealthOptions) Validate() error {
	if o.EndpointPath == "" {
		return fmt.Errorf("EndpointPath must not be empty")
	}

	if o.StatusCodeUp == 0 {
		return fmt.Errorf("StatusCodeUp must be set")
	}

	if o.StatusUp == "" {
		return fmt.Errorf("StatusUp must be set")
	}

	if o.StatusCodeDown == 0 {
		return fmt.Errorf("StatusCodeDown must be set")
	}

	if o.StatusDown == "" {
		return fmt.Errorf("StatusDown must be set")
	}

	return nil
}

// ApplyHealthOptions applies a series of HealthOption functions to a HealthOptions struct.
func ApplyHealthOptions(o *HealthOptions, opts ...HealthOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func ProvideHealthOptions(env environment.Environment) (*HealthOptions, error) {
	return options.LoadOptions[*HealthOptions](env, func(options *HealthOptions) { options.SetDefaults() })
}
