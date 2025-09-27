package config

import (
	"net/http"
	"strings"
	"time"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
	responseformat "frisboo-bank/pkg/health/enums/response_format"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
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

	// dependencies
	Logger string `mapstructure:"logger"`
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

type Option = options.OptionFn[Config]

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (Registry, error) {
	reg, err := registry.Load(
		configLoader,
		env,
		"health",
		"health",
		Default,
	)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to load health registry")
	}
	return reg, nil
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.LivenessPath, validation.Required),
		validation.Field(&c.ReadinessPath, validation.Required),
		validation.Field(&c.StatusUp, validation.Required),
		validation.Field(&c.StatusCodeUp, validation.Required),
		validation.Field(&c.StatusDown, validation.Required),
		validation.Field(&c.StatusCodeDown, validation.Required),
		validation.Field(
			&c.ResponseFormat,
			validation.Required,
			validation.By(cValidation.EnumOneOf(responseformat.ResponseFormats)),
		),
		validation.Field(&c.StartupGracePeriod, validation.Required, validation.Min(0)),
		validation.Field(&c.ShutdownDrainPeriod, validation.Required, validation.Min(0)),
		validation.Field(&c.GlobalCheckTimeout, validation.Required, validation.Min(0)),
	)
}

var Enabled = options.Option(func(c *Config, enabled bool) {
	c.Enabled = enabled
})

var LivenessPath = options.Option(func(c *Config, livenessPath string) {
	c.LivenessPath = strings.TrimSpace(livenessPath)
})

var ReadinessPath = options.Option(func(c *Config, readinessPath string) {
	c.ReadinessPath = strings.TrimSpace(readinessPath)
})

var StatusUp = options.Option(func(c *Config, statusUp string) {
	c.StatusUp = strings.TrimSpace(statusUp)
})

var StatusCodeUp = options.Option(func(c *Config, statusCodeUp int) {
	c.StatusCodeUp = statusCodeUp
})

var StatusDown = options.Option(func(c *Config, statusDown string) {
	c.StatusDown = strings.TrimSpace(statusDown)
})

var StatusCodeDown = options.Option(func(c *Config, statusCodeDown int) {
	c.StatusCodeDown = statusCodeDown
})
