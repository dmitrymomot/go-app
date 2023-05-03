package response

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmitrymomot/go-app/pkg/validator"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
)

// EncodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set(ContentTypeHeader, ContentType)

	if response == nil {
		w.WriteHeader(http.StatusCreated)
		return nil
	}

	if r, ok := response.(Responser); ok {
		return JSON(w, r)
	}

	switch r := response.(type) {
	case bool:
		return JSON(w, NewOk("", r))
	}

	return JSON(w, NewOk("", response))
}

// EncodeResponseAsIs is almost the same as EncodeResponse, but it doesn't wrap
// the response in a Response struct. This is useful for endpoints that return
// a single value, like a string or a number.
func EncodeResponseAsIs(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set(ContentTypeHeader, ContentType)

	if response == nil {
		w.WriteHeader(http.StatusCreated)
		return nil
	}

	if r, ok := response.(Responser); ok {
		return JSON(w, r)
	}

	switch r := response.(type) {
	case bool:
		return JSON(w, NewOk("", r))
	}

	return json.NewEncoder(w).Encode(response)
}

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

// EncodeError ...
func EncodeError(l logger, codeAndMessageFrom func(err error) (int, interface{})) httptransport.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		if err == nil {
			l.Log("msg", "encodeError with nil error") // nolint:errcheck
			return
		}

		code, msg := codeAndMessageFrom(err)
		if code >= http.StatusInternalServerError {
			// Log only unexpected errors
			l.Log("msg", fmt.Errorf("http transport error: %w", err)) // nolint:errcheck
		}

		var resp *Error
		switch val := msg.(type) {
		case Error:
			resp = &val
		case *Error:
			resp = val
		case *validator.ValidationError:
			resp = &Error{
				Code:       http.StatusPreconditionFailed,
				Error:      val.Err.Error(),
				Message:    "Validation error. See the validation property for more details.",
				Validation: val.Values,
			}
		default:
			resp = &Error{
				Code:    code,
				Error:   http.StatusText(code),
				Message: fmt.Sprintf("%v", msg),
			}
		}
		resp.Meta.RequestID = middleware.GetReqID(ctx)

		if err := JSON(w, resp); err != nil {
			err = fmt.Errorf("encodeError: %w", err)
			l.Log("msg", err.Error()) // nolint:errcheck
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// CodeAndMessageFrom helper
func CodeAndMessageFrom(err error) (int, interface{}) {
	if err == nil {
		return http.StatusOK, nil
	}

	if errors.Is(err, validator.ErrValidation) {
		return http.StatusPreconditionFailed, err
	}

	if errors.Is(err, jwt.ErrTokenContextMissing) {
		return http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)
	}

	if errors.Is(err, jwt.ErrTokenExpired) ||
		errors.Is(err, jwt.ErrTokenInvalid) ||
		errors.Is(err, jwt.ErrTokenMalformed) ||
		errors.Is(err, jwt.ErrTokenNotActive) ||
		errors.Is(err, jwt.ErrUnexpectedSigningMethod) {
		return http.StatusUnauthorized, err.Error()
	}

	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, err.Error()
	}

	switch err {
	case jwt.ErrTokenContextMissing:
		return http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)
	case jwt.ErrTokenExpired,
		jwt.ErrTokenInvalid,
		jwt.ErrTokenMalformed,
		jwt.ErrTokenNotActive,
		jwt.ErrUnexpectedSigningMethod:
		return http.StatusUnauthorized, err.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}
