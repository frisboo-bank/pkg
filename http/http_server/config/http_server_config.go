package config

import (
	"net"
	"strings"
	"time"

	"frisboo-bank/pkg/config"
	ginConfig "frisboo-bank/pkg/http/http_server/adapters/gin/config"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/syserrors"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*HTTPServerConfig)(nil)

type HTTPServerConfig struct {
	Type httpservertype.HttpServerType `mapstructure:"type"`

	Enabled bool `mapstructure:"enabled"`

	Host           string   `mapstructure:"host"`
	Port           string   `mapstructure:"port"`
	BasePath       string   `mapstructure:"basePath"`
	TrustedProxies []string `mapstructure:"trustedProxies"`

	Debug         bool     `mapstructure:"debug"`
	Mode          string   `mapstructure:"mode"`
	IgnoreLogUrls []string `mapstructure:"ignoreLogUrls"`

	BodyLimit             string        `mapstructure:"bodyLimit"`
	IdleTimeout           time.Duration `mapstructure:"idleTimeout"`
	MaxHeaderBytes        int           `mapstructure:"maxHeaderBytes"`
	ReadHeaderTimeout     time.Duration `mapstructure:"readHeaderTimeout"`
	ReadTimeout           time.Duration `mapstructure:"readTimeout"`
	ServerShutdownTimeout time.Duration `mapstructure:"serverShutdownTimeout"`
	WriteTimeout          time.Duration `mapstructure:"writeTimeout"`

	// adapter
	Gin ginConfig.Config `mapstructure:"gin"`

	// dependency
	Logger *loggerConfig.Config `mapstructure:"logger"`
}

func (c *HTTPServerConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func DefaultHTTPServerConfig() *HTTPServerConfig {
	return &HTTPServerConfig{
		BasePath:              "",
		BodyLimit:             "2M",
		Debug:                 false,
		Enabled:               true,
		Host:                  "0.0.0.0",
		IdleTimeout:           120 * time.Second,
		MaxHeaderBytes:        8 * 1024,
		Mode:                  "release",
		Port:                  "8080",
		ReadHeaderTimeout:     5 * time.Second,
		ReadTimeout:           30 * time.Second,
		ServerShutdownTimeout: 30 * time.Second,
		Type:                  httpservertype.HttpServerTypes.GIN,
		WriteTimeout:          30 * time.Second,
	}
}

func (c *HTTPServerConfig) Validate() error {
	var errs *multierror.Error

	if strings.TrimSpace(c.Host) == "" {
		errs = multierror.Append(errs, syserrors.CantBeEmptyError("Host"))
	}
	if c.Logger == nil {
		errs = multierror.Append(errs, syserrors.CantBeNilError("Logger"))
	} else {
		if err := c.Logger.Validate(); err != nil {
		  errs = multierror.Append(errs, syserrors.Wrap(err, "logger"))
		}
	}

	return errs.ErrorOrNil()
}
