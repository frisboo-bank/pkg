package config

import (
	"strings"
	"time"

	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
)

type Option options.OptionFn[Config]

var Host = options.OptionErr(func(c *Config, host string) error {
	host = strings.TrimSpace(host)
	if host == "" {
		return syserrors.CantBeEmptyError("Host")
	}
	c.Host = host
	return nil
})

var Port = options.OptionErr(func(c *Config, port string) error {
	port = strings.TrimSpace(port)
	if port == "" {
		return syserrors.CantBeEmptyError("Port")
	}
	c.Port = port
	return nil
})

var ServerShutdownTimeout = options.OptionErr(func(c *Config, serverShutdownTimeout time.Duration) error {
	if serverShutdownTimeout <= 0 {
		return syserrors.MustBePositiveError("ServerShutdownTimeout", serverShutdownTimeout)
	}
	c.ServerShutdownTimeout = serverShutdownTimeout
	return nil
})
