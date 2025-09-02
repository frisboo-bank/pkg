package requestid

import (
	"frisboo-bank/pkg/constants"

	"github.com/gin-gonic/gin"
)

type (
	Generator    func() string
	Handler      func(ctx *gin.Context, requestID string)
	HeaderKeyStr string
)

type Options struct {
	generator    Generator
	handler      Handler
	headerKeyStr HeaderKeyStr
	ignoredPaths []string
}

type Option func(config *Options)

var DefaultOptions = Options{
	generator:    DefaultRequestIDGenerator,
	headerKeyStr: constants.HEADER_REQUEST_ID,
	ignoredPaths: []string{},
}

func WithGenerator(generator Generator) Option {
	return func(config *Options) {
		config.generator = generator
	}
}

func WithHandler(handler Handler) Option {
	return func(config *Options) {
		config.handler = handler
	}
}

func WithHeaderKeyName(headerKeyName HeaderKeyStr) Option {
	return func(config *Options) {
		config.headerKeyStr = headerKeyName
	}
}

func WithIgnoredPaths(ignoredPaths []string) Option {
	return func(config *Options) {
		config.ignoredPaths = ignoredPaths
	}
}
