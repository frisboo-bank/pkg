package config

import (
	"net"
	"strings"
	"time"

	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/constants"
	customErrors "frisboo-bank/pkg/customerrors"
	"frisboo-bank/pkg/environment"
	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"
	"frisboo-bank/pkg/options"
)

var pError = customErrors.PrefixedError("http-server config")

type EnvConfig struct {
	BasePath              string                        `mapstructure:"basePath"`
	BodyLimit             string                        `mapstructure:"bodyLimit"`
	Development           bool                          `mapstructure:"development"`
	Host                  string                        `mapstructure:"host"`
	IdleTimeout           time.Duration                 `mapstructure:"idleTimeout"`
	IgnoreLogUrls         []string                      `mapstructure:"ignoreLogUrls"`
	MaxHeaderBytes        int                           `mapstructure:"maxHeaderBytes"`
	Port                  string                        `mapstructure:"port"`
	ReadHeaderTimeout     time.Duration                 `mapstructure:"readHeaderTimeout"`
	ReadTimeout           time.Duration                 `mapstructure:"readTimeout"`
	ServerShutdownTimeout time.Duration                 `mapstructure:"serverShutdownTimeout"`
	Type                  httpservertype.HttpServerType `mapstructure:"type"`
	WriteTimeout          time.Duration                 `mapstructure:"writeTimeout"`
}

func LoadEnvConfig(loader configContracts.ConfigLoader, env environment.Environment) (*EnvConfig, error) {
	return config.LoadConfig[EnvConfig](loader, env, "httpServer")
}

type Config struct {
	BasePath              string
	BodyLimit             string
	Development           bool
	Host                  string
	IdleTimeout           time.Duration
	IgnoreLogUrls         []string
	MaxHeaderBytes        int
	Port                  string
	ReadHeaderTimeout     time.Duration
	ReadTimeout           time.Duration
	ServerShutdownTimeout time.Duration
	WriteTimeout          time.Duration
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

var defaultConfig = &Config{
	BasePath:              "",
	BodyLimit:             constants.SERVER_BODY_LIMIT,
	Development:           false,
	Host:                  "0.0.0.0",
	IdleTimeout:           constants.SERVER_IDLE_TIMEOUT,
	MaxHeaderBytes:        constants.SERVER_MAX_HEADER_BYTES,
	Port:                  "8080",
	ReadHeaderTimeout:     constants.SERVER_READ_HEADER_TIMEOUT,
	ReadTimeout:           constants.SERVER_READ_TIMEOUT,
	ServerShutdownTimeout: constants.SERVER_SHUTDOWN_TIMEOUT,
	WriteTimeout:          constants.SERVER_WRITE_TIMEOUT,
}

func Apply() *options.OptionBuilder[Config] {
	return options.Apply(defaultConfig)
}

func FromEnvConfig(cfg *EnvConfig) *options.OptionBuilder[Config] {
	opts := Apply()

	if cfg.BasePath != "" {
		opts.With(BasePath(cfg.BasePath))
	}

	if cfg.BodyLimit != "" {
		opts.With(BodyLimit(cfg.BodyLimit))
	}

	opts.With(Development(cfg.Development))

	if cfg.Host != "" {
		opts.With(Host(cfg.Host))
	}

	if cfg.IdleTimeout != 0 {
		opts.With(IdleTimeout(cfg.IdleTimeout))
	}

	if len(cfg.IgnoreLogUrls) > 0 {
		opts.With(IgnoreLogUrls(cfg.IgnoreLogUrls))
	}

	if cfg.MaxHeaderBytes != 0 {
		opts.With(MaxHeaderBytes(cfg.MaxHeaderBytes))
	}

	if cfg.Port != "" {
		opts.With(Port(cfg.Port))
	}

	if cfg.ReadHeaderTimeout != 0 {
		opts.With(ReadHeaderTimeout(cfg.ReadHeaderTimeout))
	}

	if cfg.ReadTimeout != 0 {
		opts.With(ReadTimeout(cfg.ReadTimeout))
	}

	if cfg.ServerShutdownTimeout != 0 {
		opts.With(ServerShutdownTimeout(cfg.ServerShutdownTimeout))
	}

	if cfg.WriteTimeout != 0 {
		opts.With(WriteTimeout(cfg.WriteTimeout))
	}

	return opts
}

func BasePath(basePath string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.BasePath = strings.TrimSpace(basePath)
		return nil
	})
}

func BodyLimit(bodyLimit string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		bodyLimit = strings.TrimSpace(bodyLimit)

		if bodyLimit == "" {
			return pError.New("bodyLimit can't be empty")
		}

		cfg.BodyLimit = bodyLimit
		return nil
	})
}

func Development(development bool) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.Development = development
		return nil
	})
}

func Host(host string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		host = strings.TrimSpace(host)

		if host == "" {
			return pError.New("host can't be empty")
		}

		cfg.Host = host
		return nil
	})
}

func IdleTimeout(idleTimeout time.Duration) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if idleTimeout <= 0 {
			return pError.New("idleTimeout must be positive")
		}

		cfg.IdleTimeout = idleTimeout
		return nil
	})
}

func IgnoreLogUrls(ignoreLogUrls []string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.IgnoreLogUrls = ignoreLogUrls
		return nil
	})
}

func MaxHeaderBytes(maxHeaderBytes int) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if maxHeaderBytes <= 0 {
			return pError.New("maxHeaderBytes must be positive")
		}

		cfg.MaxHeaderBytes = maxHeaderBytes
		return nil
	})
}

func Port(port string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		port = strings.TrimSpace(port)

		if port == "" {
			return pError.New("port can't be empty")
		}

		cfg.Port = port
		return nil
	})
}

func ReadHeaderTimeout(readHeaderTimeout time.Duration) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if readHeaderTimeout <= 0 {
			return pError.New("readHeaderTimeout must be positive")
		}

		cfg.ReadHeaderTimeout = readHeaderTimeout
		return nil
	})
}

func ReadTimeout(readTimeout time.Duration) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if readTimeout <= 0 {
			return pError.New("readTimeout must be positive")
		}

		cfg.ReadTimeout = readTimeout
		return nil
	})
}

func ServerShutdownTimeout(serverShutdownTimeout time.Duration) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if serverShutdownTimeout <= 0 {
			return pError.New("serverShutdownTimeout must be positive")
		}

		cfg.ServerShutdownTimeout = serverShutdownTimeout
		return nil
	})
}

func WriteTimeout(writeTimeout time.Duration) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if writeTimeout <= 0 {
			return pError.New("writeTimeout must be positive")
		}

		cfg.WriteTimeout = writeTimeout
		return nil
	})
}
