package options

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"
	httpservertype "frisboo-bank/pkg/http/http_server/options/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
	optionsContracts "frisboo-bank/pkg/options/contracts"
)

type HttpServerOptions struct {
	Type                  httpservertype.HttpServerType `mapstructure:"type"`
	BasePath              string                        `mapstructure:"basePath"`
	Development           bool                          `mapstructure:"development"`
	Host                  string                        `mapstructure:"host"`
	Port                  string                        `mapstructure:"port"`
	IgnoreLogUrls         []string                      `mapstructure:"ignoreLogUrls"`
	BodyLimit             string                        `mapstructure:"bodyLimit"`
	IdleTimeout           time.Duration                 `mapstructure:"idleTimeout"`
	MaxHeaderBytes        int                           `mapstructure:"maxHeaderBytes"`
	ReadHeaderTimeout     time.Duration                 `mapstructure:"readHeaderTimeout"`
	ReadTimeout           time.Duration                 `mapstructure:"readTimeout"`
	ServerShutdownTimeout time.Duration                 `mapstructure:"serverShutdownTimeout"`
	WriteTimeout          time.Duration                 `mapstructure:"writeTimeout"`
	Logger                loggerContracts.Logger
}

var defaultOptions = HttpServerOptions{
	Type:                  httpservertype.HttpServerTypes.GIN,
	BasePath:              "",
	Development:           false,
	Host:                  "0.0.0.0",
	Port:                  "8080",
	IgnoreLogUrls:         []string{},
	BodyLimit:             constants.SERVER_BODY_LIMIT,
	IdleTimeout:           constants.SERVER_IDLE_TIMEOUT,
	MaxHeaderBytes:        constants.SERVER_MAX_HEADER_BYTES,
	ReadHeaderTimeout:     constants.SERVER_READ_HEADER_TIMEOUT,
	ReadTimeout:           constants.SERVER_READ_TIMEOUT,
	ServerShutdownTimeout: constants.SERVER_SHUTDOWN_TIMEOUT,
	WriteTimeout:          constants.SERVER_WRITE_TIMEOUT,
}

var _ optionsContracts.Options = (*HttpServerOptions)(nil)

func (o *HttpServerOptions) Clone() optionsContracts.Options {
	return options.CloneOptions(o)
}

// SetDefaults implements contracts.Options.
func (o *HttpServerOptions) SetDefaults() {
	if !o.Type.IsValid() {
		o.Type = defaultOptions.Type
	}

	if o.BasePath == "" {
		o.BasePath = defaultOptions.BasePath
	}

	if o.Host == "" {
		o.Host = defaultOptions.Host
	}

	if o.Port == "" {
		o.Port = defaultOptions.Port
	}

	if o.IgnoreLogUrls == nil {
		o.IgnoreLogUrls = defaultOptions.IgnoreLogUrls
	}

	if o.BodyLimit == "" {
		o.BodyLimit = defaultOptions.BodyLimit
	}

	if o.IdleTimeout <= 0 {
		o.IdleTimeout = defaultOptions.IdleTimeout
	}

	if o.MaxHeaderBytes <= 0 {
		o.MaxHeaderBytes = defaultOptions.MaxHeaderBytes
	}

	if o.ReadHeaderTimeout <= 0 {
		o.ReadHeaderTimeout = defaultOptions.ReadHeaderTimeout
	}

	if o.ReadTimeout <= 0 {
		o.ReadTimeout = defaultOptions.ReadTimeout
	}

	if o.ServerShutdownTimeout <= 0 {
		o.ServerShutdownTimeout = defaultOptions.ServerShutdownTimeout
	}

	if o.WriteTimeout <= 0 {
		o.WriteTimeout = defaultOptions.WriteTimeout
	}
}

func (o *HttpServerOptions) Validate() error {
	if o.Host == "" {
		return fmt.Errorf("host must not be empty")
	}

	if o.Port == "" {
		return fmt.Errorf("port must not be empty")
	}
	return nil
}

func (o *HttpServerOptions) Address() string {
	return net.JoinHostPort(o.Host, o.Port)
}

func (o *HttpServerOptions) FullAddress() (string, error) {
	path, err := url.JoinPath(o.Address(), o.BasePath)
	if err != nil {
		return "", err
	}

	return path, nil
}

// ApplyHttpServerOptions applies a series of HttpServerOption functions to a HttpServerOptions struct.
func ApplyHttpServerOptions(o *HttpServerOptions, opts ...HttpServerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func ProvideHttpServerOptions(env environment.Environment) (*HttpServerOptions, error) {
	return options.LoadOptions[*HttpServerOptions](env, func(options *HttpServerOptions) { options.SetDefaults() })
}
