package jsonrpc

import (
	"encoding/json"
	"net/http"
)

// NewJSONRPCResponse creates new JSON-RPC 2.0 response
func NewJSONRPCResponse(id interface{}, result interface{}) *JSONRPCResponse {
	return &JSONRPCResponse{
		Version: "2.0",
		ID:      id,
		Result:  result,
	}
}

// NewJSONRPCError creates new JSON-RPC 2.0 error
func NewJSONRPCError(id interface{}, code int, message string, data interface{}) *JSONRPCResponse {
	return &JSONRPCResponse{
		Version: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// NewJSONRPCRequest creates new JSON-RPC 2.0 request
func NewJSONRPCRequest(method string, params interface{}, id interface{}) *JSONRPCRequest {
	return &JSONRPCRequest{
		Version: "2.0",
		Method:  method,
		Params:  params,
		ID:      id,
	}
}

// NewJSONRPCNotification creates new JSON-RPC 2.0 notification
func NewJSONRPCNotification(method string, params interface{}) *JSONRPCRequest {
	return &JSONRPCRequest{
		Version: "2.0",
		Method:  method,
		Params:  params,
	}
}

// ParseRequest parses JSON-RPC 2.0 request from JSON string
func ParseRequest(data []byte) (*JSONRPCRequest, error) {
	var req JSONRPCRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// ParseResponse parses JSON-RPC 2.0 response from JSON string
func ParseResponse(data []byte) (*JSONRPCResponse, error) {
	var res JSONRPCResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ParseHTTPRequest parses JSON-RPC 2.0 request from HTTP request
func ParseHTTPRequest(req *http.Request) (*JSONRPCRequest, error) {
	var jsonReq JSONRPCRequest
	if err := json.NewDecoder(req.Body).Decode(&jsonReq); err != nil {
		return nil, err
	}
	return &jsonReq, nil
}
