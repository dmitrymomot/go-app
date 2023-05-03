package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Predefined http encoder content type
const (
	ContentTypeHeader = "Content-Type"
	ContentType       = "application/json; charset=utf-8"
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
func (resp *defaultResponder) JSON(w http.ResponseWriter, response Responser, headersKV ...string) error {
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

	// set default headers
	w.Header().Set(ContentTypeHeader, ContentType)
	// set custom headers
	for i := 0; i < len(headersKV); i += 2 {
		w.Header().Set(headersKV[i], headersKV[i+1])
	}

	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response.GetPayload()); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}
	return nil
}
