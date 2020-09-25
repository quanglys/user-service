package http

import (
	"context"
	"encoding/json"
	"net/http"
	"user-service/src/service/transport"
)

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == http.ErrHandlerTimeout {
		return
	}

	e, ok := err.(transport.Error)
	if !ok {
		panic("encodeError invalid error")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(codeToHTTPStatus(e.Code))
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": transport.ErrorResponse{
			Code: e.Code,
			Msg:  e.Error(),
		},
	})
}

func codeToHTTPStatus(code transport.ResponseCode) int {
	status := http.StatusInternalServerError
	switch code {
	case transport.ErrorCodeInvalidParameter:
		status = http.StatusBadRequest
	case transport.ErrorCodePermissionDenied:
		status = http.StatusForbidden
	case transport.ErrorCodeNotFound:
		status = http.StatusNotFound
	case transport.ErrorCodeEmpty:
		status = http.StatusNoContent
	case transport.ErrorCodeNotImplemented:
		status = http.StatusNotImplemented
	case transport.ErrorCodeUnauthorized:
		status = http.StatusUnauthorized
	default:
		status = http.StatusInternalServerError
	}
	return status
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(*transport.Error); ok && e != nil {
		encodeErrorResponse(ctx, e, w)
		return nil
	} else if r, ok := response.(transport.ResponseCode); ok {
		w.WriteHeader(codeToHTTPStatus(r))
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

