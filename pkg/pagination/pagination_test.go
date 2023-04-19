package pagination_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmitrymomot/go-server/pkg/pagination"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		params      map[string]interface{}
		expected    *pagination.Pagination
	}{
		{
			name:        "With JSON body",
			method:      http.MethodPost,
			contentType: "application/json",
			params: map[string]interface{}{
				"limit":  20,
				"offset": 40,
			},
			expected: &pagination.Pagination{
				Limit:  20,
				Offset: 40,
				Page:   3,
			},
		},
		{
			name:        "With query params",
			method:      http.MethodGet,
			contentType: "",
			params: map[string]interface{}{
				"limit":  "30",
				"offset": "60",
				"page":   "0",
			},
			expected: &pagination.Pagination{
				Limit:  30,
				Offset: 60,
				Page:   3,
			},
		},
		{
			name:        "empty params",
			method:      http.MethodGet,
			contentType: "text/plain",
			params:      map[string]interface{}{},
			expected: &pagination.Pagination{
				Limit:  10,
				Offset: 0,
				Page:   1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data io.Reader
			var req *http.Request

			if tt.method == http.MethodPost && tt.contentType == "application/json" {
				data = strings.NewReader(buildJSON(tt.params))
				req = httptest.NewRequest(tt.method, "/", data)
				if tt.contentType != "" {
					req.Header.Set("Content-Type", tt.contentType)
				}
			} else {
				req = httptest.NewRequest(tt.method, "/?"+buildQueryString(tt.params), data)
				if tt.contentType != "" {
					req.Header.Set("Content-Type", tt.contentType)
				}
			}

			p := pagination.Parse(req)
			assert.Equal(t, tt.expected, p)
		})
	}
}

func TestGetOffset(t *testing.T) {
	tests := []struct {
		name     string
		p        *pagination.Pagination
		expected int
	}{
		{
			name: "With offset",
			p: &pagination.Pagination{
				Limit:  10,
				Offset: 30,
				Page:   2,
			},
			expected: 30,
		},
		{
			name: "Without offset",
			p: &pagination.Pagination{
				Limit:  20,
				Offset: 0,
				Page:   3,
			},
			expected: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := tt.p.GetOffset()
			assert.Equal(t, tt.expected, offset)
		})
	}
}

func TestGetPages(t *testing.T) {
	tests := []struct {
		name     string
		p        *pagination.Pagination
		total    int
		expected int
	}{
		{
			name: "Exact multiple of limit",
			p: &pagination.Pagination{
				Limit:  10,
				Offset: 0,
				Page:   1,
			},
			total:    20,
			expected: 2,
		},
		{
			name: "Not exact multiple of limit",
			p: &pagination.Pagination{
				Limit:  15,
				Offset: 30,
				Page:   3,
			},
			total:    37,
			expected: 3,
		},
		{
			name: "Zero total or limit",
			p: &pagination.Pagination{
				Limit:  5,
				Offset: 0,
				Page:   1,
			},
			total:    0,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages := tt.p.GetPages(tt.total)
			assert.Equal(t, tt.expected, pages)
		})
	}
}

func buildQueryString(params map[string]interface{}) string {
	var parts []string
	for k, v := range params {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(parts, "&")
}

func buildJSON(params map[string]interface{}) string {
	json, _ := json.Marshal(params)
	return string(json)
}
