package config

import (
	"strings"
	"time"

	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"

	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	loggerConfig "frisboo-bank/pkg/logger/config"
)

type HTTPServerOption = options.OptionFn[HTTPServerConfig]

var Type = options.OptionErr(func(c *HTTPServerConfig, sType httpservertype.HttpServerType) error {
	if err := validation.EnumOneOf("Type", sType, httpservertype.HttpServerTypes); err != nil {
		return err
	}
	c.Type = sType
	return nil
})

var Enabled = options.Option(func(c *HTTPServerConfig, enabled bool) {
	c.Enabled = enabled
})

var BasePath = options.Option(func(c *HTTPServerConfig, basePath string) {
	c.BasePath = strings.TrimSpace(basePath)
})

var BodyLimit = options.OptionErr(func(c *HTTPServerConfig, bodyLimit string) error {
	bodyLimit = strings.TrimSpace(bodyLimit)
	if bodyLimit == "" {
		return syserrors.CantBeEmptyError("BodyLimit")
	}
	c.BodyLimit = bodyLimit
	return nil
})

var Debug = options.Option(func(c *HTTPServerConfig, debug bool) {
	c.Debug = debug
})

var Host = options.OptionErr(func(c *HTTPServerConfig, host string) error {
	host = strings.TrimSpace(host)
	if host == "" {
		return syserrors.CantBeEmptyError("Host")
	}
	c.Host = host
	return nil
})

var IdleTimeout = options.OptionErr(func(c *HTTPServerConfig, idleTimeout time.Duration) error {
	if idleTimeout <= 0 {
		return syserrors.MustBePositiveError("IdleTimeout", idleTimeout)
	}
	c.IdleTimeout = idleTimeout
	return nil
})

var IgnoreLogUrls = options.Option(func(c *HTTPServerConfig, ignoreLogUrls []string) {
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

var AppendIgnoreLogUrls = options.VarOption(func(c *HTTPServerConfig, ignoreLogUrls ...string) {
	for _, p := range ignoreLogUrls {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		c.IgnoreLogUrls = append(c.IgnoreLogUrls, p)
	}
})

var Logger = options.OptionErr(func(c *HTTPServerConfig, logger *loggerConfig.Config) error {
	if err := validation.NotNil("Logger", logger); err != nil {
		return err
	}
	c.Logger = logger
	return nil
})

var LoggerOptions = options.VarOptionErr(func(c *HTTPServerConfig, opts ...loggerConfig.Option) error {
	l := c.Logger
	if err := options.Apply(l, opts...); err != nil {
		return err
	}
	if err := validation.NotNil("Logger", l); err != nil {
		return err
	}
	c.Logger = l
	return nil
})

var MaxHeaderBytes = options.OptionErr(func(c *HTTPServerConfig, maxHeaderBytes int) error {
	if maxHeaderBytes <= 0 {
		return syserrors.MustBePositiveError("MaxHeaderBytes", maxHeaderBytes)
	}
	c.MaxHeaderBytes = maxHeaderBytes
	return nil
})

var Mode = options.Option(func(c *HTTPServerConfig, mode string) {
	c.Mode = strings.TrimSpace(mode)
})

var Port = options.OptionErr(func(c *HTTPServerConfig, port string) error {
	port = strings.TrimSpace(port)
	if port == "" {
		return syserrors.CantBeEmptyError("Port")
	}
	c.Port = port
	return nil
})

var ReadHeaderTimeout = options.OptionErr(func(c *HTTPServerConfig, readHeaderTimeout time.Duration) error {
	if readHeaderTimeout <= 0 {
		return syserrors.MustBePositiveError("ReadHeaderTimeout", readHeaderTimeout)
	}
	c.ReadHeaderTimeout = readHeaderTimeout
	return nil
})

var ReadTimeout = options.OptionErr(func(c *HTTPServerConfig, readTimeout time.Duration) error {
	if readTimeout <= 0 {
		return syserrors.MustBePositiveError("ReadTimeout", readTimeout)
	}
	c.ReadTimeout = readTimeout
	return nil
})

var ServerShutdownTimeout = options.OptionErr(func(c *HTTPServerConfig, serverShutdownTimeout time.Duration) error {
	if serverShutdownTimeout <= 0 {
		return syserrors.MustBePositiveError("ServerShutdownTimeout", serverShutdownTimeout)
	}
	c.ServerShutdownTimeout = serverShutdownTimeout
	return nil
})

var TrustedProxies = options.Option(func(c *HTTPServerConfig, proxies []string) {
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

var AppendTrustedProxies = options.VarOption(func(c *HTTPServerConfig, proxies ...string) {
	for _, p := range proxies {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		c.TrustedProxies = append(c.TrustedProxies, p)
	}
})

var WriteTimeout = options.OptionErr(func(c *HTTPServerConfig, writeTimeout time.Duration) error {
	if writeTimeout <= 0 {
		return syserrors.MustBePositiveError("WriteTimeout", writeTimeout)
	}
	c.WriteTimeout = writeTimeout
	return nil
})
