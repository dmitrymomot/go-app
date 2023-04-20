package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// global default responder
var dr Responder = NewDefaultResponder()

// SetDefaultResponder sets default responder.
func SetDefaultResponder(responder Responder) {
	dr = responder
}

// defaultResponder is a default implementation of Responder.
type defaultResponder struct{}

// NewDefaultResponder creates new default responder.
func NewDefaultResponder() Responder {
	return &defaultResponder{}
}

// NewResponder creates new Responder.
func NewResponder() Responder {
	return &defaultResponder{}
}

// JSON HTTP response
func (resp *defaultResponder) JSON(w http.ResponseWriter, response Responser) error {
	if response == nil {
		return fmt.Errorf("response is nil")
	}
	code := response.GetCode()
	if response.GetCode() == 0 {
		code = http.StatusOK
	}

	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response.GetPayload()); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}
	return nil
}
