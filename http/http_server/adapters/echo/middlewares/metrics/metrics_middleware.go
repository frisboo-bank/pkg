package metrics

import (
	"frisboo-bank/pkg/http/http_server/adapters/echo/middlewares/metrics/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/validation"

	"github.com/labstack/echo/v4"
)

type httpMetrics struct{}

func HTTPMetrics(logger loggerContracts.Logger, opts ...config.Option) (echo.MiddlewareFunc, error) {
	validation.AssertNotNil("logger", logger)

	cfg, err := config.New(opts...)
	if err != nil {
		return nil, err
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return nil
		}
	}, nil
}
