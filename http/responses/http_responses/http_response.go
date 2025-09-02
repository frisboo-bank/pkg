package httpresponses

import "frisboo-bank/pkg/syserrors"

type HTTPResponse struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

func NewHTTPErrorResponse(err syserrors.E, withStack bool) *HTTPResponse {
	if err == nil {
		return nil
	}

	resp := &HTTPResponse{
		Status: 0,
	}

	return resp
}
