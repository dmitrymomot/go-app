package response

// New creates new response
func New(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewError creates new error response
func NewError(code int, err error, message string, validation map[string][]string) *Error {
	return &Error{
		Code:       code,
		Error:      err.Error(),
		Message:    message,
		Validation: validation,
	}
}

// NewList creates new list response
func NewList(items interface{}, total int, pagination *Pagination) *List {
	return &List{
		Items:      items,
		Total:      total,
		Pagination: pagination,
	}
}

// NewPagination creates new pagination
func NewPagination(limit, offset, page, pages int) *Pagination {
	return &Pagination{
		Limit:  limit,
		Offset: offset,
		Page:   page,
		Pages:  pages,
	}
}
