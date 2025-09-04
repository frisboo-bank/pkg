package httperrorresponse

import (
	"net/http"
	"testing"

	"frisboo-bank/pkg/syserrors"

	"github.com/stretchr/testify/assert"
)

func TestBuilderDefault(t *testing.T) {
	resp := NewHTTPErrorResponseBuilder().
		Build()

	assert.Equal(t, DefaultStatus, resp.Status)
	assert.Equal(t, DefaultCode, resp.Code)
	assert.Equal(t, DefaultMessage, resp.Message)
	assert.Nil(t, resp.Details)
	assert.Empty(t, resp.Stack)
}

func TestBuilderSetter(t *testing.T) {
	builder := NewHTTPErrorResponseBuilder().
		SetStatus(http.StatusNotFound).
		SetCode("not_found").
		SetMessage("resource missing").
		AddDetail("request_id", "customer-request-id").
		SetStack("dummy stack").
		WithStack(true)

	resp := builder.Build()

	assert.Equal(t, http.StatusNotFound, resp.Status)
	assert.Equal(t, "not_found", resp.Code)
	assert.Equal(t, "resource missing", resp.Message)
	assert.Len(t, resp.Details, 1)
	assert.Equal(t, "customer-request-id", resp.Details["request_id"])
	assert.Equal(t, "dummy stack", resp.Stack)
}

func TestDetailsToggle(t *testing.T) {
	builder := NewHTTPErrorResponseBuilder().
		SetDetails(map[string]any{"resource_id": "12345"})

	resp := builder.
		WithDetails(true).
		Build()

	assert.Len(t, resp.Details, 1)
	assert.Equal(t, "12345", resp.Details["resource_id"])

	resp = builder.
		WithDetails(false).
		Build()

	assert.Nil(t, resp.Details)
}

func TestStackToggle(t *testing.T) {
	builder := NewHTTPErrorResponseBuilder().
		SetStack("stack trace here")

	resp := builder.
		WithStack(true).
		Build()

	assert.Equal(t, "stack trace here", resp.Stack)

	resp = builder.
		WithStack(false).
		Build()

	assert.Empty(t, resp.Stack)
}

func TestFromError(t *testing.T) {
	err := syserrors.WithDetails(
		syserrors.New("db connection failed"),
		"code", "unavailable",
		"status", http.StatusServiceUnavailable,
		"endpoint", "/user/id",
		"id", "1234",
	)

	builder := FromError(err).
		WithStack(true)

	resp := builder.
		Build()

	assert.Equal(t, http.StatusServiceUnavailable, resp.Status)
	assert.Equal(t, "unavailable", resp.Code)
	assert.Equal(t, "db connection failed", resp.Message)
	assert.Len(t, resp.Details, 4)
	assert.Equal(t, "unavailable", resp.Details["code"])
	assert.Equal(t, http.StatusServiceUnavailable, resp.Details["status"])
	assert.Equal(t, "/user/id", resp.Details["endpoint"])
	assert.Equal(t, "1234", resp.Details["id"])
	assert.NotEmpty(t, resp.Stack)

	resp = builder.
		WithStack(false).
		Build()

	assert.Empty(t, resp.Stack)
}
