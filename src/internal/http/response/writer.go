package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)

	if err := enc.Encode(payload); err != nil {
		http.Error(w, `{"error":{"code":"internal_error","message":"internal server error"}}`, http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, status int, code string, message string, details []FieldError) {
	resp := ErrorResponse{
		Error: Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	WriteJSON(w, status, resp)
}

func WriteAppError(w http.ResponseWriter, err error) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		details := []FieldError{}
		for field, msg := range ve.Fields {
			details = append(details, FieldError{
				Field:   field,
				Message: msg,
			})
		}

		WriteError(w, http.StatusBadRequest, "validation_error", "invalid request", details)
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		WriteError(w, http.StatusNotFound, "not_found", "resource not found", nil)
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			WriteError(w, http.StatusConflict, "conflict", "resource already exists", nil)
		case "23503", "23502", "22P02":
			WriteError(w, http.StatusBadRequest, "validation_error", "invalid request", nil)
		default:
			WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error", nil)
		}
		return
	}

	WriteError(w, http.StatusInternalServerError, "internar_error", fmt.Sprintf("internal server error: %s", err), nil)
	return
}
