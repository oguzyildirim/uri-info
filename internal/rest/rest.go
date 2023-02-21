package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Oguzyildirim/url-info/internal"
	"go.opentelemetry.io/otel/trace"
)

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Error string `json:"error"`
}

func renderErrorResponse(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	resp := ErrorResponse{Error: msg}
	status := http.StatusInternalServerError

	var ierr *internal.Error
	if !errors.As(err, &ierr) {
		resp.Error = "internal error"
	} else {
		switch ierr.Code() {
		case internal.ErrorCodeNotFound:
			status = http.StatusNotFound
		case internal.ErrorCodeInvalidArgument:
			status = http.StatusBadRequest
		}
	}

	if err != nil {
		_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("UrlTracer").Start(ctx, "rest.renderErrorResponse")
		defer span.End()

		span.RecordError(err)
	}

	renderResponse(w, resp, status)
}

func renderResponse(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	content, err := json.Marshal(res)
	if err != nil {
		// error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if _, err = w.Write(content); err != nil {
		//  error
	}
}
