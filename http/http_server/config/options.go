package config

import (
	"strings"
	"time"

	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

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
