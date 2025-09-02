package httperrorresponse

import (
	"frisboo-bank/pkg/syserrors"
	"maps"
	"net/http"
)

const (
	DefaultCode    = "unknown"
	DefaultStatus  = http.StatusInternalServerError
	DefaultMessage = "something went wrong"
)

// OpenAPI schema (conceptual):
// type: object
// required: [status, message]
type HTTPErrorResponse struct {
	Status  int            `json:"status"            example:"404"`
	Code    string         `json:"code,omitempty"    example:"not_found"`
	Message string         `json:"message"           example:"resource not found"`
	Details map[string]any `json:"details,omitempty"`
	Stack   string         `json:"stack,omitempty"`
}

type HTTPErrorResponseBuilder interface {
	SetStatus(status int) HTTPErrorResponseBuilder
	SetCode(code string) HTTPErrorResponseBuilder
	SetMessage(message string) HTTPErrorResponseBuilder
	SetDetails(details map[string]any) HTTPErrorResponseBuilder
	AddDetail(key string, value any) HTTPErrorResponseBuilder
	SetStack(stack string) HTTPErrorResponseBuilder

	WithStack(withStack bool) HTTPErrorResponseBuilder
	WithDetails(withDetails bool) HTTPErrorResponseBuilder

	Build() *HTTPErrorResponse
}

var _ HTTPErrorResponseBuilder = (*httpErrorResponseBuilder)(nil)

type httpErrorResponseBuilder struct {
	status      int
	code        string
	message     string
	details     map[string]any
	stack       string
	withDetails bool
	withStack   bool
}

func NewHTTPErrorResponseBuilder() HTTPErrorResponseBuilder {
	return &httpErrorResponseBuilder{
		status:      DefaultStatus,
		code:        DefaultCode,
		message:     DefaultMessage,
		withDetails: true,
		withStack:   false,
	}
}

// FromError seeds the builder from a syserrors.E.
func FromError(err error) HTTPErrorResponseBuilder {
	builder := NewHTTPErrorResponseBuilder()
	if err == nil {
		return builder
	}

	details := syserrors.AllDetails(err)

	return NewHTTPErrorResponseBuilder().
		SetStatus(inferStatusFromDetails(details)).
		SetCode(syserrors.FormatCode(err)).
		SetMessage(err.Error()).
		SetDetails(details).
		SetStack(syserrors.FormatStack(err))
}

func (h *httpErrorResponseBuilder) SetCode(code string) HTTPErrorResponseBuilder {
	h.code = code
	return h
}

func (h *httpErrorResponseBuilder) SetDetails(details map[string]any) HTTPErrorResponseBuilder {
	if len(details) == 0 {
		h.details = nil
	}

	out := make(map[string]any, len(details))
	maps.Copy(out, details)
	h.details = out
	return h
}

func (h *httpErrorResponseBuilder) AddDetail(key string, value any) HTTPErrorResponseBuilder {
	if len(h.details) == 0 {
		h.details = make(map[string]any, 1)
	}
	h.details[key] = value
	return h
}

func (h *httpErrorResponseBuilder) SetMessage(message string) HTTPErrorResponseBuilder {
	h.message = message
	return h
}

func (h *httpErrorResponseBuilder) SetStack(stack string) HTTPErrorResponseBuilder {
	h.stack = stack
	return h
}

func (h *httpErrorResponseBuilder) SetStatus(status int) HTTPErrorResponseBuilder {
	if status < 100 || status >= 600 {
		status = DefaultStatus
	}
	h.status = status
	return h
}

func (h *httpErrorResponseBuilder) WithDetails(withDetails bool) HTTPErrorResponseBuilder {
	h.withDetails = withDetails
	return h
}

func (h *httpErrorResponseBuilder) WithStack(withStack bool) HTTPErrorResponseBuilder {
	h.withStack = withStack
	return h
}

func (h *httpErrorResponseBuilder) Build() *HTTPErrorResponse {
	resp := &HTTPErrorResponse{
		Status:  h.status,
		Code:    h.code,
		Message: h.message,
	}

	if h.withDetails && len(h.details) > 0 {
		resp.Details = h.details
	}

	if h.withStack && h.stack != "" {
		resp.Stack = h.stack
	}

	return resp
}

func inferStatusFromDetails(details map[string]any) int {
	if details == nil {
		return DefaultStatus
	}

	if status, ok := details["status"].(int); ok {
		return status
	}

	return DefaultStatus
}
