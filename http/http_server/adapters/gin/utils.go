package gin

import (
	"frisboo-bank/pkg/syserrors"

	"github.com/gin-gonic/gin"
)

func ToMiddlewareType(middleware any) (gin.HandlerFunc, error) {
	switch mw := middleware.(type) {
	case gin.HandlerFunc:
		return mw, nil
	case func(*gin.Context):
		return mw, nil
	}

	return nil, syserrors.New("invalid middleware type, must be gin.HandlerFunc or func(*gin.Context)")
}

func ToMiddlewaresType(middlewares ...any) ([]gin.HandlerFunc, error) {
	mws := make([]gin.HandlerFunc, len(middlewares))

	for i, m := range middlewares {
		mw, err := ToMiddlewareType(m)
		if err != nil {
			return nil, syserrors.Wrapf(err, "error while converting middleware %d", i)
		}

		mws[i] = mw
	}

	return mws, nil
}
