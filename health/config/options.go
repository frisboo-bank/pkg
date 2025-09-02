package config

import (
	"strings"

	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
)

type Option = options.OptionFn[Config]

var Enabled = options.Option(func(c *Config, enabled bool) {
	c.Enabled = enabled
})

var LivenessPath = options.OptionErr(func(c *Config, livenessPath string) error {
	livenessPath = strings.TrimSpace(livenessPath)
	if livenessPath == "" {
		return syserrors.CantBeEmptyError("LivenessPath")
	}
	c.LivenessPath = livenessPath
	return nil
})

var ReadinessPath = options.OptionErr(func(c *Config, readinessPath string) error {
	readinessPath = strings.TrimSpace(readinessPath)
	if readinessPath == "" {
		return syserrors.CantBeEmptyError("ReadinessPath")
	}
	c.ReadinessPath = readinessPath
	return nil
})

var StatusUp = options.OptionErr(func(c *Config, statusUp string) error {
	statusUp = strings.TrimSpace(statusUp)
	if statusUp == "" {
		return syserrors.CantBeEmptyError("StatusUp")
	}
	c.StatusUp = statusUp
	return nil
})

var StatusCodeUp = options.OptionErr(func(c *Config, statusCodeUp int) error {
	if statusCodeUp <= 0 {
		return syserrors.MustBePositiveError("StatusCodeUp", statusCodeUp)
	}
	c.StatusCodeUp = statusCodeUp
	return nil
})

var StatusDown = options.OptionErr(func(c *Config, statusDown string) error {
	statusDown = strings.TrimSpace(statusDown)
	if statusDown == "" {
		return syserrors.CantBeEmptyError("StatusDown")
	}
	c.StatusDown = statusDown
	return nil
})

var StatusCodeDown = options.OptionErr(func(c *Config, statusCodeDown int) error {
	if statusCodeDown <= 0 {
		return syserrors.New("statusCodeDown must be positive")
	}
	c.StatusCodeDown = statusCodeDown
	return nil
})
