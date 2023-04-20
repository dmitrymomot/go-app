package jsonrpc

type (
	// JSON-RPC 2.0 response
	// https://www.jsonrpc.org/specification#response_object
	JSONRPCResponse struct {
		// JSON-RPC version
		Version string `json:"jsonrpc" example:"2.0"`
		// Result of the method call
		Result interface{} `json:"result,omitempty" example:"{}"`
		// Error object
		Error *JSONRPCError `json:"error,omitempty"`
		// Request ID
		ID interface{} `json:"id,omitempty" example:"1"`
	}

	// JSON-RPC 2.0 error
	// https://www.jsonrpc.org/specification#error_object
	JSONRPCError struct {
		// Error code
		Code int `json:"code" example:"-32603"`
		// Error message
		Message string `json:"message" example:"Internal error"`
		// Error data
		Data interface{} `json:"data,omitempty" example:"{}"`
	}

	// JSON-RPC 2.0 request
	// https://www.jsonrpc.org/specification#request_object
	JSONRPCRequest struct {
		// JSON-RPC version
		Version string `json:"jsonrpc" example:"2.0"`
		// Method name
		Method string `json:"method" example:"method"`
		// Request parameters
		Params interface{} `json:"params,omitempty" example:"{}"`
		// Request ID
		ID interface{} `json:"id,omitempty" example:"1"`
	}
)
