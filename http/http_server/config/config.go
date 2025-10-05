package config

import (
	"net"
	"strings"
	"time"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
	echoConfig "frisboo-bank/pkg/http/http_server/adapters/echo/config"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	cValidation "frisboo-bank/pkg/validation"

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
	GzipLevel             int                           `mapstructure:"gzipLevel"`

	// adapters
	Echo echoConfig.Config `mapstructure:"echo"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (c Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func Default() Config {
	return Config{
		Enabled:               true,
		Debug:                 false,
		Mode:                  "release",
		Host:                  "127.0.0.1",
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
		GzipLevel:             5,
	}
}

type (
	Option   = options.OptionFn[Config]
	Registry = registry.Registry[Config]
)

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (Registry, error) {
	reg, err := registry.Load(
		configLoader,
		env,
		"httpServers",
		"httpServer",
		Default,
	)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to load http-server registry")
	}
	return reg, nil
}

func (c *Config) Validate() error {
	if err := validation.ValidateStruct(c,
		validation.Field(&c.Type, validation.Required, validation.By(cValidation.EnumOneOf(httpservertype.HttpServerTypes))),
		validation.Field(&c.Host, validation.Required, validationIs.Host),
		validation.Field(&c.Port, validation.Required, validationIs.Port),
		validation.Field(&c.BasePath, validation.Required),
		validation.Field(&c.Mode, validation.Required, validation.In("debug", "release", "test")),
		validation.Field(&c.IgnoreLogUrls),
		validation.Field(&c.TrustedProxies),
		validation.Field(&c.MaxHeaderBytes, validation.Required, validation.Min(1024)),
		validation.Field(&c.BodyLimit, validation.Required),
		validation.Field(&c.IdleTimeout, validation.Required, validation.Min(1*time.Second)),
		validation.Field(&c.ReadHeaderTimeout, validation.Required, validation.Min(1*time.Second)),
		validation.Field(&c.ReadTimeout, validation.Required, validation.Min(1*time.Second)),
		validation.Field(&c.ServerShutdownTimeout, validation.Required, validation.Min(1*time.Second)),
		validation.Field(&c.WriteTimeout, validation.Required, validation.Min(1*time.Second)),
		validation.Field(&c.GzipLevel, validation.Required, validation.Min(-1), validation.Max(9)),
	); err != nil {
		return err
	}

	switch c.Type {
	case httpservertype.HttpServerTypes.ECHO:
		if err := validation.Validate(&c.Echo, validation.Required); err != nil {
			return err
		}
		return c.Echo.Validate()
	}

	return nil
}

var Type = options.Option(func(c *Config, sType httpservertype.HttpServerType) {
	c.Type = sType
})

var Enabled = options.Option(func(c *Config, enabled bool) {
	c.Enabled = enabled
})

var BasePath = options.Option(func(c *Config, basePath string) {
	c.BasePath = strings.TrimSpace(basePath)
})

var BodyLimit = options.Option(func(c *Config, bodyLimit string) {
	c.BodyLimit = strings.TrimSpace(bodyLimit)
})

var Debug = options.Option(func(c *Config, debug bool) {
	c.Debug = debug
})

var Host = options.Option(func(c *Config, host string) {
	c.Host = strings.TrimSpace(host)
})

var IdleTimeout = options.Option(func(c *Config, idleTimeout time.Duration) {
	c.IdleTimeout = idleTimeout
})

var IgnoreLogUrls = options.Option(func(c *Config, ignoreLogUrls []string) {
	out := make([]string, 0, len(ignoreLogUrls))
	for _, p := range ignoreLogUrls {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	c.IgnoreLogUrls = out
})

var AppendIgnoreLogUrls = options.VarOption(func(c *Config, ignoreLogUrls ...string) {
	for _, p := range ignoreLogUrls {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		c.IgnoreLogUrls = append(c.IgnoreLogUrls, p)
	}
})

var MaxHeaderBytes = options.Option(func(c *Config, maxHeaderBytes int) {
	c.MaxHeaderBytes = maxHeaderBytes
})

var Mode = options.Option(func(c *Config, mode string) {
	c.Mode = strings.TrimSpace(mode)
})

var Port = options.Option(func(c *Config, port string) {
	c.Port = strings.TrimSpace(port)
})

var ReadHeaderTimeout = options.Option(func(c *Config, readHeaderTimeout time.Duration) {
	c.ReadHeaderTimeout = readHeaderTimeout
})

var ReadTimeout = options.Option(func(c *Config, readTimeout time.Duration) {
	c.ReadTimeout = readTimeout
})

var ServerShutdownTimeout = options.Option(func(c *Config, serverShutdownTimeout time.Duration) {
	c.ServerShutdownTimeout = serverShutdownTimeout
})

var TrustedProxies = options.Option(func(c *Config, proxies []string) {
	out := make([]string, 0, len(proxies))
	for _, p := range proxies {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	c.TrustedProxies = out
})

var AppendTrustedProxies = options.VarOption(func(c *Config, proxies ...string) {
	for _, p := range proxies {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		c.TrustedProxies = append(c.TrustedProxies, p)
	}
})

var WriteTimeout = options.Option(func(c *Config, writeTimeout time.Duration) {
	c.WriteTimeout = writeTimeout
})

var GzipLevel = options.Option(func(c *Config, gzipLevel int) {
	c.GzipLevel = gzipLevel
})

var EchoConfig = options.Option(func(c *Config, echoConfig echoConfig.Config) {
	c.Echo = echoConfig
})

var Logger = options.Option(func(c *Config, logger string) {
	c.Logger = logger
})
