package options

import (
	"time"

	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"
	httpservertype "frisboo-bank/pkg/http/http_server/options/enums/http_server_type"
)

var (
	Type                  = httpservertype.HttpServerTypes.GIN
	BasePath              = ""
	Host                  = "0.0.0.0"
	Port                  = "8080"
	BodyLimit             = constants.SERVER_BODY_LIMIT
	IdleTimeout           = constants.SERVER_IDLE_TIMEOUT
	MaxHeaderBytes        = constants.SERVER_MAX_HEADER_BYTES
	ReadHeaderTimeout     = constants.SERVER_READ_HEADER_TIMEOUT
	ReadTimeout           = constants.SERVER_READ_TIMEOUT
	ServerShutdownTimeout = constants.SERVER_SHUTDOWN_TIMEOUT
	WriteTimeout          = constants.SERVER_WRITE_TIMEOUT
)

type HTTPServerOptions struct {
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
}

func ProvideHTTPServerOptions(
	loader configContracts.ConfigLoader,
	env environment.Environment,
) (*HTTPServerOptions, error) {
	return config.LoadOptions[HTTPServerOptions](loader, env)
}
