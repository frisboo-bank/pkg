package gin

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ToMiddlewareType(middleware any) (gin.HandlerFunc, error) {
	switch mw := middleware.(type) {
	case gin.HandlerFunc:
		return mw, nil
	case func(*gin.Context):
		return mw, nil
	}

	return nil, fmt.Errorf("invalid middleware type, must be gin.HandlerFunc or func(*gin.Context)")
}

func ToMiddlewaresType(middlewares ...any) ([]gin.HandlerFunc, error) {
	mws := make([]gin.HandlerFunc, len(middlewares))

	for i, m := range middlewares {
		mw, err := ToMiddlewareType(m)
		if err != nil {
			return nil, fmt.Errorf("error while add middleware %d: %w", i, err)
		}

		mws[i] = mw
	}

	return mws, nil
}
