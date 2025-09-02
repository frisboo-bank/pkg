package httpresponses

import "frisboo-bank/pkg/syserrors"

type HTTPResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func NewHTTPErrorResponse(err syserrors.E, withStack bool) *HTTPResponse {
	if err == nil {
		return nil
	}

	resp := &HTTPResponse{
		Status:  0,
		Code:    "",
		Message: err.Error(),
	}

	details := syserrors.AllDetails(err)
	if len(details) > 0 {
		var cDetails map[string]any
		copy(cDetails, syserrors.AllDetails(err))
		resp.Details = cDetails
	}

	if withStack {
		resp.Stack = syserrors.FormatStack(err)
	}

	return resp
}
