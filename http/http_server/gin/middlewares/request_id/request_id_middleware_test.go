package requestid

import (
	"frisboo-bank/pkg/constants"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddlewre(t *testing.T) {
	customHeaderKeyName := "my-customer-header-key"
	customRequestID := "my-custom-request-id"
	alreadySetRequestID := "my-already-set-Request-id"

	tests := []struct {
		name              string
		middleware        gin.HandlerFunc
		incomingRequestId string
		expectation       func(w *httptest.ResponseRecorder)
	}{
		{
			name:       "Generates a new Request ID",
			middleware: NewRequestIDMiddleware(),
			expectation: func(w *httptest.ResponseRecorder) {
				requestID := w.Header().Get(constants.HEADER_REQUEST_ID)

				assert.Equal(t, 36, len(requestID), "The request_id length should be 36")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			router.Use(tt.middleware)

			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "Response ok")
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)

			if tt.incomingRequestId != "" {
				req.Header.Add(constants.HEADER_REQUEST_ID, tt.incomingRequestId)
			}

			router.ServeHTTP(w, req)

			tt.expectation(w)
		})
	}
}
