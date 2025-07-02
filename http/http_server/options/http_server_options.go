package options

import (
	"net"
	"net/url"
	"time"

	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/reflection/typemapper"

	"github.com/stoewer/go-strcase"
)

type (
	HttpServerType string
)

const (
	TypeGin = HttpServerType("gin")
)

const (
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderXFrameOptions                 = "X-Frame-Options"
	HeaderCacheControl                  = "Cache-Control"
	HeaderPragma                        = "Pragma"
	HeaderExpires                       = "Expires"
	HeaderStrictTransportSecurity       = "Strict-Transport-Security"
	HeaderContentSecurityPolicy         = "Content-Security-Policy"
	HeaderContentTypeOptions            = "X-Content-Type-Options"
	HeaderXSSProtection                 = "X-XSS-Protection"
	HeaderReferrerPolicy                = "Referrer-Policy"
	HeaderPermissionsPolicy             = "Permissions-Policy"
	HeaderServer                        = "Server"
	HeaderCrossOriginResourcePolicy     = "Cross-Origin-Resource-Policy"
)

type HttpServerOptions struct {
	Type                  HttpServerType `mapstructure:"type"`
	BasePath              string         `mapstructure:"basePath"`
	Development           bool           `mapstructure:"development"`
	Host                  string         `mapstructure:"host"`
	Port                  string         `mapstructure:"port"`
	IgnoreLogUrls         []string       `mapstructure:"ignoreLogUrls"`
	BodyLimit             string         `mapstructure:"bodyLimit"`
	IdleTimeout           time.Duration  `mapstructure:"idleTimeout"`
	MaxHeaderBytes        int            `mapstructure:"maxHeaderBytes"`
	ReadHeaderTimeout     time.Duration  `mapstructure:"readHeaderTimeout"`
	ReadTimeout           time.Duration  `mapstructure:"readTimeout"`
	ServerShutdownTimeout time.Duration  `mapstructure:"serverShutdownTimeout"`
	WriteTimeout          time.Duration  `mapstructure:"writeTimeout"`
	Logger                loggerContracts.Logger
}

type HttpServerOption = func(options *HttpServerOptions)

var DefaultHttpServerOptions = HttpServerOptions{
	Type:                  TypeGin,
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

func WithLogger(logger loggerContracts.Logger) HttpServerOption {
	return func(options *HttpServerOptions) {
		options.Logger = logger
	}
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

var optionName = strcase.LowerCamelCase(typemapper.GetGenericTypeNameByT[HttpServerOptions]())

func ProvideHttpServerOptions(env environment.Environment) (*HttpServerOptions, error) {
	opt, err := config.BindConfigKey[*HttpServerOptions](optionName, env)
	if err != nil {
		return nil, err
	}

	return setDefaultOptions(opt), nil
}

func setDefaultOptions(options *HttpServerOptions) *HttpServerOptions {
	if options.Type == "" {
		options.Type = DefaultHttpServerOptions.Type
	}

	if options.BasePath == "" {
		options.BasePath = DefaultHttpServerOptions.BasePath
	}

	if options.Host == "" {
		options.Host = DefaultHttpServerOptions.Host
	}

	if options.Port == "" {
		options.Port = DefaultHttpServerOptions.Port
	}

	if options.IgnoreLogUrls == nil {
		options.IgnoreLogUrls = DefaultHttpServerOptions.IgnoreLogUrls
	}

	if options.BodyLimit == "" {
		options.BodyLimit = DefaultHttpServerOptions.BodyLimit
	}

	if options.IdleTimeout <= 0 {
		options.IdleTimeout = DefaultHttpServerOptions.IdleTimeout
	}

	if options.MaxHeaderBytes <= 0 {
		options.MaxHeaderBytes = DefaultHttpServerOptions.MaxHeaderBytes
	}

	if options.ReadHeaderTimeout <= 0 {
		options.ReadHeaderTimeout = DefaultHttpServerOptions.ReadHeaderTimeout
	}

	if options.ReadTimeout <= 0 {
		options.ReadTimeout = DefaultHttpServerOptions.ReadTimeout
	}

	if options.ServerShutdownTimeout <= 0 {
		options.ServerShutdownTimeout = DefaultHttpServerOptions.ServerShutdownTimeout
	}
	if options.WriteTimeout <= 0 {
		options.WriteTimeout = DefaultHttpServerOptions.WriteTimeout
	}

	return options
}
