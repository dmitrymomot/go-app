package response

import "net/http"

type (
	// Response interface
	Responser interface {
		// GetCode returns HTTP code of response
		GetCode() int
		// GetPayload returns payload of response
		GetPayload() interface{}
	}

	// Responder is an interface for response.
	// Respond writes response to http.ResponseWriter.
	Responder interface {
		// JSON writes JSON response to http.ResponseWriter.
		// headersKV is a list of headers key-value pairs.
		// E.g. "Content-Type", "application/json", "X-Request-ID", "123"
		JSON(w http.ResponseWriter, response Responser, headersKV ...string) error
	}

	// Response is a struct for response
	Response struct {
		// Code is a code of response
		Code int `json:"code" example:"200"`
		// Message is a message of response
		Message string `json:"message" example:"OK"`
		// Data is a data of response
		Data interface{} `json:"data" example:"{}"`
		// Meta is a meta of response
		Meta *Meta `json:"meta,omitempty"`
	}

	// Error is a struct for error response
	Error struct {
		// Code is a code of error
		Code int `json:"code" example:"400"`
		// Error is a string error representation
		Error string `json:"error" example:"Bad Request"`
		// Message is a message of error
		Message string `json:"message" example:"Validation error"`
		// Validation is a validation of error
		Validation map[string][]string `json:"validation,omitempty" example:"{'email': ['Email is required']}"`
		// Meta is a meta of response
		Meta *Meta `json:"meta,omitempty"`
	}

	// List represents a list of response data
	List struct {
		// Items is a list of items
		Items interface{} `json:"items" example:"[]"`
		// Total is a total count of items
		Total int `json:"total" example:"0"`
		// Pagination is a pagination of list
		Pagination *Pagination `json:"pagination,omitempty"`
	}

	// Pagination represents a pagination of list
	Pagination struct {
		// Limit is a limit of items per page
		Limit int `json:"limit"`
		// Offset is a offset of items per page
		Offset int `json:"offset"`
		// Page is a current page
		Page int `json:"page"`
		// Pages is a total count of pages
		Pages int `json:"pages"`
	}

	// Meta represents a meta of response
	Meta struct {
		// Title is a title of response
		Title string `json:"title" example:"Title"`
		// Description is a description of response
		Description string `json:"description" example:"Description"`
		// Version is a version of response
		Version string `json:"version" example:"1.0.0"`
		// RequestID is a request ID of response
		RequestID string `json:"request_id"`
	}
)

// GetCode returns HTTP code of response
func (r *Response) GetCode() int {
	return r.Code
}

// GetPayload returns payload of response
func (r *Response) GetPayload() interface{} {
	return r
}

// GetCode returns HTTP code of response
func (e *Error) GetCode() int {
	return e.Code
}

// GetPayload returns payload of response
func (e *Error) GetPayload() interface{} {
	return e
}
