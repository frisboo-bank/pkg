package config

import (
	"net"
	"strings"
	"time"

	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	ginConfig "frisboo-bank/pkg/http/http_server/adapters/gin/config"
	loggerConfig "frisboo-bank/pkg/logger/config"

	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	Type httpservertype.HttpServerType `mapstructure:"type"`

	Enabled bool `mapstructure:"enabled"`

	Host           string   `mapstructure:"host"`
	Port           string   `mapstructure:"port"`
	BasePath       string   `mapstructure:"basePath"`
	TrustedProxies []string `mapstructure:"trustedProxies"`

	Debug bool   `mapstructure:"debug"`
	Mode  string `mapstructure:"mode"`

	BodyLimit             string        `mapstructure:"bodyLimit"`
	IdleTimeout           time.Duration `mapstructure:"idleTimeout"`
	MaxHeaderBytes        int           `mapstructure:"maxHeaderBytes"`
	ReadHeaderTimeout     time.Duration `mapstructure:"readHeaderTimeout"`
	ReadTimeout           time.Duration `mapstructure:"readTimeout"`
	ServerShutdownTimeout time.Duration `mapstructure:"serverShutdownTimeout"`
	WriteTimeout          time.Duration `mapstructure:"writeTimeout"`

	Gin ginConfig.Config `mapstructure:"gin"`

	Logger        loggerConfig.Config `mapstructure:"logger"`
	IgnoreLogUrls []string            `mapstructure:"ignoreLogUrls"`
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func Default() *Config {
	loggerCfg := loggerConfig.Default()
	loggerCfg.Prefix = "http-server"

	return &Config{
		Mode:                  "release",
		BasePath:              "",
		BodyLimit:             "2M",
		Debug:                 false,
		Enabled:               true,
		Host:                  "0.0.0.",
		IdleTimeout:           120 * time.Second,
		Logger:                *loggerCfg,
		MaxHeaderBytes:        8 * 1024,
		Port:                  "8080",
		ReadHeaderTimeout:     5 * time.Second,
		ReadTimeout:           30 * time.Second,
		ServerShutdownTimeout: 30 * time.Second,
		Type:                  httpservertype.HttpServerTypes.GIN,
		WriteTimeout:          30 * time.Second,
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	if strings.TrimSpace(c.Host) == "" {
		errs = multierror.Append(errs, syserrors.CantBeEmptyError("Host"))
	}
	if err := c.Logger.Validate(); err != nil {
		errs = multierror.Append(errs, syserrors.Wrap(err, "logger"))
	}

	return errs.ErrorOrNil()
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (*Config, error) {
	cfg := Default()
	if err := loader.LoadByKey("http-server", env, cfg); err != nil {
		return nil, err
	}
	return options.New(func() *Config { return cfg }, opts...)
}
