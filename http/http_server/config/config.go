package config

import (
	"net"
	"time"

	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"

	cValidation "frisboo-bank/pkg/validation"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"

	ginConfig "frisboo-bank/pkg/http/http_server/adapters/gin/config"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	validationIs "github.com/go-ozzo/ozzo-validation/v4/is"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Enabled               bool                          `mapstructure:"enabled"`
	Type                  httpservertype.HttpServerType `mapstructure:"type"`
	Debug                 bool                          `mapstructure:"debug"`
	Mode                  string                        `mapstructure:"mode"`
	Host                  string                        `mapstructure:"host"`
	Port                  string                        `mapstructure:"port"`
	BasePath              string                        `mapstructure:"basePath"`
	IgnoreLogUrls         []string                      `mapstructure:"ignoreLogUrls"`
	TrustedProxies        []string                      `mapstructure:"trustedProxies"`
	MaxHeaderBytes        int                           `mapstructure:"maxHeaderBytes"`
	BodyLimit             string                        `mapstructure:"bodyLimit"`
	IdleTimeout           time.Duration                 `mapstructure:"idleTimeout"`
	ReadHeaderTimeout     time.Duration                 `mapstructure:"readHeaderTimeout"`
	ReadTimeout           time.Duration                 `mapstructure:"readTimeout"`
	ServerShutdownTimeout time.Duration                 `mapstructure:"serverShutdownTimeout"`
	WriteTimeout          time.Duration                 `mapstructure:"writeTimeout"`

	// adapters
	Gin *ginConfig.Config `mapstructure:"gin"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func Default() Config {
	return Config{
		Enabled:               true,
		Type:                  httpservertype.HttpServerTypes.GIN,
		Debug:                 false,
		Mode:                  "release",
		Host:                  "0.0.0.0",
		Port:                  "8080",
		BasePath:              "/",
		TrustedProxies:        nil,
		IgnoreLogUrls:         nil,
		MaxHeaderBytes:        8 * 1024,
		BodyLimit:             "2M",
		IdleTimeout:           120 * time.Second,
		ReadHeaderTimeout:     5 * time.Second,
		ReadTimeout:           30 * time.Second,
		ServerShutdownTimeout: 30 * time.Second,
		WriteTimeout:          30 * time.Second,
	}
}

func (c *Config) Validate() error {
	if err := validation.ValidateStruct(c,
		validation.Field(&c.Type, validation.Required, validation.By(cValidation.EnumOneOf(httpservertype.HttpServerTypes))),
		validation.Field(&c.Host, validation.Required, validationIs.Host),
		validation.Field(&c.Port, validation.Required, validationIs.Port),
		validation.Field(&c.BasePath, validation.Required),
		validation.Field(&c.IgnoreLogUrls),
		validation.Field(&c.MaxHeaderBytes, validation.Required, validation.Min(0)),
		validation.Field(&c.BodyLimit, validation.Required),
		validation.Field(&c.IdleTimeout, validation.Required, validation.Min(0)),
		validation.Field(&c.ReadHeaderTimeout, validation.Required, validation.Min(0)),
		validation.Field(&c.ReadTimeout, validation.Required, validation.Min(0)),
		validation.Field(&c.ServerShutdownTimeout, validation.Required, validation.Min(0)),
		validation.Field(&c.WriteTimeout, validation.Required, validation.Min(0)),
	); err != nil {
		return err
	}

	switch c.Type {
	case httpservertype.HttpServerTypes.GIN:
		if err := validation.Validate(&c.Gin, validation.Required); err != nil {
			return err
		}
		return c.Gin.Validate()
	}

	return nil
}

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (*Registry, error) {
	return registry.Load(
		configLoader,
		env,
		"httpServers",
		"httpServer",
		Default,
	)
}
