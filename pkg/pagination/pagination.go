package pagination

import (
	"net/http"
	"strings"

	"github.com/dmitrymomot/binder"
)

type (
	// Pagination represents a pagination of list
	Pagination struct {
		// Limit is a limit of items per page
		Limit int `json:"limit" query:"limit" example:"10"`
		// Offset is a offset of items per page
		Offset int `json:"offset" query:"offset" example:"0"`
		// Page is a current page
		Page int `json:"page" query:"page" example:"1"`
	}
)

// Parse parses pagination from request
func Parse(r *http.Request) *Pagination {
	p := &Pagination{}

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		// Parse pagination from JSON body if Content-Type is application/json
		_ = binder.BindJSON(r, p) // nolint: errcheck
	} else {
		// Parse pagination from query string
		_ = binder.BindQuery(r, p) // nolint: errcheck
	}

	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Page == 0 {
		if p.Offset == 0 {
			p.Page = 1
		} else {
			p.Page = (p.Offset / p.Limit) + 1
		}
	}

	return p
}

// GetLimit returns limit of pagination
func (p *Pagination) GetLimit() int {
	if p.Limit > 0 {
		return p.Limit
	}
	return 10
}

// GetOffset returns offset of pagination
func (p *Pagination) GetOffset() int {
	if p.Offset > 0 {
		return p.Offset
	}
	return (p.Page - 1) * p.Limit
}

// GetPage returns current page
func (p *Pagination) GetPage() int {
	if p.Page > 0 {
		return p.Page
	}
	return (p.Offset / p.Limit) + 1
}

// GetPages returns total count of pages
func (p *Pagination) GetPages(total int) int {
	if total == 0 || p.Limit == 0 {
		return 1
	}
	if total%p.Limit == 0 {
		return total / p.Limit
	}
	return total/p.Limit + 1
}
