package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dmitrymomot/go-server/pkg/response"
	"github.com/dmitrymomot/go-utils"
)

// Testing middleware. Helps to test any HTTP error.
// Pass must_err query parameter with code you want get
// E.g.: /login?must_err=403
func Testing() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errCodeStr := r.URL.Query().Get("must_err")
			if len(errCodeStr) == 0 {
				// Call next middleware if must_err is not set
				next.ServeHTTP(w, r)
				return
			}

			errCode, err := strconv.Atoi(errCodeStr)
			if err != nil {
				// Return the error with status code 400
				sendErrorResponse(w, r, http.StatusBadRequest, err)
				return
			}

			if !isValidErrorCode(errCode) {
				// Return invalid error code error
				sendErrorResponse(w, r, http.StatusBadRequest, errors.New("Invalid error code"))
				return
			}

			// Return error with needed status code
			if utils.IsJsonRequest(r) {
				response.JSON(w, response.NewError(
					errCode,
					errors.New(http.StatusText(errCode)),
					http.StatusText(errCode),
					nil,
				))
			} else {
				http.Error(w, http.StatusText(errCode), errCode)
			}
		})
	}
}

// Helper function to check if an error code is valid
func isValidErrorCode(errCode int) bool {
	return errCode >= 400 && errCode < 600
}

// Helper function to send an error response
func sendErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	if utils.IsJsonRequest(r) {
		response.JSON(w, response.NewError(
			statusCode,
			err,
			http.StatusText(statusCode),
			nil,
		))
	} else {
		http.Error(w, err.Error(), statusCode)
	}
}
