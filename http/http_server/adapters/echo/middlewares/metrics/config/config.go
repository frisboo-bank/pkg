package config

import (
	"frisboo-bank/pkg/options"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	skipper         middleware.Skipper
	metricsProvider metric.MeterProvider
	name            string
}

type Option = options.OptionFn[Config]

func Default() Config {
	return Config{
		skipper:         middleware.DefaultSkipper,
		metricsProvider: otel.GetMeterProvider(),
	}
}

func New(opts ...Option) (Config, error) {
	var zero Config

	base := Default()
	if err := options.Apply(&base, opts...); err != nil {
		return zero, err
	}

	return base, nil
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.skipper, validation.Required),
		validation.Field(&c.metricsProvider, validation.Required),
		validation.Field(&c.name, validation.Required),
	)
}

var Skipper = options.Option(func(c *Config, skipper middleware.Skipper) {
	c.skipper = skipper
})

var MetricsProvider = options.Option(func(c *Config, metricsProvider metric.MeterProvider) {
	c.metricsProvider = metricsProvider
})

var Name = options.Option(func(c *Config, name string) {
	c.name = name
})
