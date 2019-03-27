package service

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_ParsePagingParams(t *testing.T) {

	type testCase struct {
		name     string
		request  func(t *testing.T) *http.Request
		validate func(t *testing.T, req *http.Request)
	}

	tests := []testCase{
		{
			name: "happy path",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?page_number=1&page_size=100&sort_asc=true&sort_field=foo", nil)
				assert.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				pagingParams, err := ParsePagingParams(req, []string{"foo", "bar"})
				assert.Nil(t, err)

				assert.Equal(t, int32(1), *pagingParams.PageNumber)
				assert.Equal(t, int32(100), *pagingParams.PageSize)
				assert.Equal(t, true, *pagingParams.SortAsc)
				assert.Equal(t, "foo", *pagingParams.SortField)
			},
		},
		{
			name: "sort field not whitelisted",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?page_number=1&page_size=100&sort_asc=true&sort_field=foo", nil)
				assert.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				_, err := ParsePagingParams(req, []string{"bar"})
				assert.NotNil(t, err)
			},
		},
		{
			name: "sort field not included",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?page_number=1&page_size=100&sort_asc=true", nil)
				assert.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				pagingParams, err := ParsePagingParams(req, []string{"bar"})
				assert.Nil(t, err)

				assert.Equal(t, int32(1), *pagingParams.PageNumber)
				assert.Equal(t, int32(100), *pagingParams.PageSize)
				assert.Equal(t, true, *pagingParams.SortAsc)
				assert.Nil(t, pagingParams.SortField)
			},
		},
		{
			name: "page number not an integer",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?page_number=foo&page_size=100&sort_asc=true&sort_field=foo", nil)
				assert.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				_, err := ParsePagingParams(req, []string{"foo"})
				assert.NotNil(t, err)
			},
		},
		{
			name: "sort asc boolean parse error",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?page_number=1&page_size=100&sort_asc=nope&sort_field=foo", nil)
				assert.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				_, err := ParsePagingParams(req, []string{"foo"})
				assert.NotNil(t, err)
			},
		},
		{
			name: "sort asc missing, does not error",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?page_number=1&page_size=100&sort_field=foo", nil)
				assert.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				pagingParams, err := ParsePagingParams(req, []string{"foo", "bar"})
				assert.Nil(t, err)

				assert.Equal(t, int32(1), *pagingParams.PageNumber)
				assert.Equal(t, int32(100), *pagingParams.PageSize)
				assert.Equal(t, "foo", *pagingParams.SortField)
				assert.Nil(t, pagingParams.SortAsc)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.validate(t, tc.request(t))
		})
	}
}
