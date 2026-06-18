package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nilotpaul/stock-market-agg-svc/model"
)

const requestIDKey = "request_id"

type ApiHandlerFunc func(http.ResponseWriter, *http.Request) error

func NewApiError(status int, msg string, data any) *model.ApiError {
	return &model.ApiError{
		Status:  status,
		Message: msg,
		Data:    data,
	}
}

func NewBadRequestError(msg string, data any) *model.ApiError {
	return &model.ApiError{
		Status:  http.StatusBadRequest,
		Message: msg,
		Data:    data,
	}
}

func NewInternalServerError(msg string, data any) *model.ApiError {
	return &model.ApiError{
		Status:  http.StatusInternalServerError,
		Message: msg,
		Data:    data,
	}
}

func NewNotFoundError(msg string, data any) *model.ApiError {
	return &model.ApiError{
		Status:  http.StatusNotFound,
		Message: msg,
		Data:    data,
	}
}

// requestID in logs are turned off for debugging in development
// but should be included in production.
func handler(fn ApiHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			start = time.Now()
			reqID = uuid.New().String()
		)

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		r = r.WithContext(ctx)

		err := fn(w, r)
		slog.Info(
			"http_request",
			// "request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"took_ms", time.Since(start).Milliseconds(),
		)

		if err == nil {
			return
		}

		slog.Error(
			"http_error",
			// "request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"err", err.Error(),
			"took_ms", time.Since(start).Milliseconds(),
		)

		var apiErr *model.ApiError
		if errors.As(err, &apiErr) {
			writeJSON(w, apiErr.Status, apiErr)
		} else {
			writeJSON(w, http.StatusInternalServerError, model.ApiError{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
			})
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
