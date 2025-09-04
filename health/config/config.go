package config

import (
	"net/http"
	"strings"
	"time"

	"frisboo-bank/pkg/config"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/environment"
	responseformat "frisboo-bank/pkg/health/enums/response_format"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	Enabled bool `mapstructure:"enabled"`

	LivenessPath  string `mapstructure:"livenessPath"`
	ReadinessPath string `mapstructure:"readinessPath"`

	StatusUp       string `mapstructure:"statusUp"`
	StatusCodeUp   int    `mapstructure:"statusCodeUp"`
	StatusDown     string `mapstructure:"statusDown"`
	StatusCodeDown int    `mapstructure:"statusCodeDown"`

	AdditionalHeaders map[string]string             `mapstructure:"additionalHeaders"`
	ResponseFormat    responseformat.ResponseFormat `mapstructure:"responseFormat"`

	StartupGracePeriod  time.Duration `mapstructure:"startupGracePeriod"`  // During this window readiness may be forced "Down" or "Degraded".
	ShutdownDrainPeriod time.Duration `mapstructure:"shutdownDrainPeriod"` // Time to keep reporting not-ready so traffic drains.
	GlobalCheckTimeout  time.Duration `mapstructure:"globalCheckTimeout"`  // Upper bound across all dependency checks (0 = disabled).

	// dependency
	Logger *loggerConfig.Config `mapstructure:"logger"`
}

func Default() *Config {
	return &Config{
		Enabled:             true,
		LivenessPath:        "/healthz",
		ReadinessPath:       "/readyz",
		StatusUp:            "UP",
		StatusCodeUp:        http.StatusOK,
		StatusDown:          "DOWN",
		StatusCodeDown:      http.StatusServiceUnavailable,
		AdditionalHeaders:   map[string]string{},
		ResponseFormat:      responseformat.ResponseFormats.JSON,
		StartupGracePeriod:  15 * time.Second,
		ShutdownDrainPeriod: 5 * time.Second,
		GlobalCheckTimeout:  5 * time.Second,
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	if !c.Enabled {
		return nil
	}

	if strings.TrimSpace(c.LivenessPath) == "" {
		errs = multierror.Append(errs, syserrors.CantBeEmptyError("LivenessPath"))
	}
	if strings.TrimSpace(c.ReadinessPath) == "" {
		errs = multierror.Append(errs, syserrors.CantBeEmptyError("ReadinessPath"))
	}

	return errs.ErrorOrNil()
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (*Config, error) {
	cfg := Default()
	if err := loader.LoadByKey("health", env, cfg); err != nil {
		return nil, err
	}
	return options.New(func() *Config { return cfg }, opts...)
}
