package config

import (
	"net/http"
	"time"

	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	responseformat "frisboo-bank/pkg/health/enums/response_format"
	loggerConfig "frisboo-bank/pkg/logger/config"

	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

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

func Default() Config {
	return Config{
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
	return validation.ValidateStruct(c,
		validation.Field(&c.LivenessPath, validation.Required),
		validation.Field(&c.ReadinessPath, validation.Required),
		validation.Field(&c.StatusUp, validation.Required),
		validation.Field(&c.StatusCodeUp, validation.Required),
		validation.Field(&c.StatusDown, validation.Required),
		validation.Field(&c.StatusCodeDown, validation.Required),
		validation.Field(&c.ResponseFormat, validation.Required, validation.By(cValidation.EnumOneOf(responseformat.ResponseFormats))),
		validation.Field(&c.StartupGracePeriod, validation.Required, validation.Min(0)),
		validation.Field(&c.ShutdownDrainPeriod, validation.Required, validation.Min(0)),
		validation.Field(&c.GlobalCheckTimeout, validation.Required, validation.Min(0)),
	)
}

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (*Registry, error) {
	return registry.Load(
		configLoader,
		env,
		"health",
		"health",
		Default,
	)
}
