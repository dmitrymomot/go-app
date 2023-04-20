package response

// New creates new response
func New(code int, message string, data interface{}, meta *Meta) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
		Meta:    meta,
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

// NewMeta creates new meta
func NewMeta(title, description, version string) *Meta {
	return &Meta{
		Title:       title,
		Description: description,
		Version:     version,
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

// paginator is an interface for pagination
type paginator interface {
	// GetLimit returns limit
	GetLimit() int
	// GetOffset returns offset
	GetOffset() int
	// GetPage returns page
	GetPage() int
	// GetPages returns pages
	GetPages() int
}

// NewPaginationFromInterface creates new pagination from interface
func NewPaginationFromInterface(source paginator) *Pagination {
	return &Pagination{
		Limit:  source.GetLimit(),
		Offset: source.GetOffset(),
		Page:   source.GetPage(),
		Pages:  source.GetPages(),
	}
}
