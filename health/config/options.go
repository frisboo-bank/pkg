package config

import (
	"strings"

	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

var Enabled = options.Option(func(c *Config, enabled bool) {
	c.Enabled = enabled
})

var LivenessPath = options.Option(func(c *Config, livenessPath string) {
	c.LivenessPath = strings.TrimSpace(livenessPath)
})

var ReadinessPath = options.Option(func(c *Config, readinessPath string) {
	c.ReadinessPath = strings.TrimSpace(readinessPath)
})

var StatusUp = options.Option(func(c *Config, statusUp string) {
	c.StatusUp = strings.TrimSpace(statusUp)
})

var StatusCodeUp = options.Option(func(c *Config, statusCodeUp int) {
	c.StatusCodeUp = statusCodeUp
})

var StatusDown = options.Option(func(c *Config, statusDown string) {
	c.StatusDown = strings.TrimSpace(statusDown)
})

var StatusCodeDown = options.Option(func(c *Config, statusCodeDown int) {
	c.StatusCodeDown = statusCodeDown
})
