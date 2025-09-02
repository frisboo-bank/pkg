package requestid

import (
	"strings"

	"frisboo-bank/pkg/syserrors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewRequestIDMiddleware(options ...Option) gin.HandlerFunc {
	opts := DefaultOptions

	for _, f := range options {
		f(&opts)
	}

	return newRequestIDMiddleware(&opts)
}

func newRequestIDMiddleware(options *Options) gin.HandlerFunc {
	headerKeyStr := string(options.headerKeyStr)

	return func(ctx *gin.Context) {
		requestID := strings.TrimSpace(ctx.Request.Header.Get(headerKeyStr))

		if requestID == "" {
			requestID = options.generator()
			ctx.Request.Header.Add(headerKeyStr, requestID)
		}

		if options.handler != nil {
			options.handler(ctx, requestID)
		}

		ctx.Header(headerKeyStr, requestID)

		ctx.Next()
	}
}

func DefaultRequestIDGenerator() string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(syserrors.Newf("request-id: failed to generate id with error: %w", err))
	}

	return id.String()
}
