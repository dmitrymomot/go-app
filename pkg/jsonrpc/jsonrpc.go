package jsonrpc

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
