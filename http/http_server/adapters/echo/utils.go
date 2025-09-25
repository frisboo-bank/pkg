package echo

import (
	"frisboo-bank/pkg/syserrors"

	echoVendor "github.com/labstack/echo/v4"
)

func ToMiddlewareType(middleware any) (echoVendor.MiddlewareFunc, error) {
	switch mw := middleware.(type) {
	case echoVendor.MiddlewareFunc:
		return mw, nil
	}

	return nil, syserrors.New("invalid middleware type, must be echo.MiddlewareFunc")
}

func ToMiddlewaresType(middlewares ...any) ([]echoVendor.MiddlewareFunc, error) {
	mws := make([]echoVendor.MiddlewareFunc, len(middlewares))

	for i, m := range middlewares {
		mw, err := ToMiddlewareType(m)
		if err != nil {
			return nil, syserrors.Wrapf(err, "error while converting middleware %d", i)
		}

		mws[i] = mw
	}

	return mws, nil
}
