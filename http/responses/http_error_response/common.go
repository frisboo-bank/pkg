package httperrorresponse

import "net/http"

type httpCommonErrorResponse struct {
	Status int
	Code   string
}

var (
	ErrorNotFound = httpCommonErrorResponse{
		Status: http.StatusNotFound,
		Code:   "not_found",
	}
	ErrorUnauthorized = httpCommonErrorResponse{
		Status: http.StatusUnauthorized,
		Code:   "unauthorized",
	}
	ErrorForbidden = httpCommonErrorResponse{
		Status: http.StatusForbidden,
		Code:   "forbidden",
	}
	ErrorConflict = httpCommonErrorResponse{
		Status: http.StatusConflict,
		Code:   "conflict",
	}
	ErrorServiceUnavailable = httpCommonErrorResponse{
		Status: http.StatusServiceUnavailable,
		Code:   "service_unavailable",
	}
	ErrorInvalidInput = httpCommonErrorResponse{
		Status: http.StatusBadRequest,
		Code:   "invalid_input",
	}
	ErrorInternalServerError = httpCommonErrorResponse{
		Status: http.StatusInternalServerError,
		Code:   "internal_error",
	}
	ErrorTimeout = httpCommonErrorResponse{
		Status: http.StatusGatewayTimeout,
		Code:   "timeout",
	}
	ErrorDeadline = httpCommonErrorResponse{
		Status: http.StatusGatewayTimeout,
		Code:   "deadline",
	}
	ErrorAlreadyExists = httpCommonErrorResponse{
		Status: http.StatusConflict,
		Code:   "already_exists",
	}
	ErrorRateLimited = httpCommonErrorResponse{
		Status: http.StatusTooManyRequests,
		Code:   "rate_limited",
	}
)

func ErrorNotFoundResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorNotFound, err)
}

func ErrorUnauthorizedResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorUnauthorized, err)
}

func ErrorForbiddenResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorForbidden, err)
}

func ErrorConflictResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorConflict, err)
}

func ErrorUnavailableResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorServiceUnavailable, err)
}

func ErrorInvalidInputResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorInvalidInput, err)
}

func ErrorInternalServerResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorInternalServerError, err)
}

func ErrorTimeoutResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorTimeout, err)
}

func ErrorDeadlineResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorDeadline, err)
}

func ErrorAlreadyExistsResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorAlreadyExists, err)
}

func ErrorRateLimitedResponse(err error) HTTPErrorResponseBuilder {
	return fromHTTPCommonErrorResponse(ErrorRateLimited, err)
}

func fromHTTPCommonErrorResponse(commonError httpCommonErrorResponse, err error) HTTPErrorResponseBuilder {
	return FromError(err).
		SetStatus(commonError.Status).
		SetCode(commonError.Code)
}
