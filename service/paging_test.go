package service

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnit_SortValid(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T)
	}{
		"base path": {
			validate: func(t *testing.T) {
				sort := Sort{
					Field:     "foo",
					Direction: "asc",
				}

				whiteList := []string{"bar", "hello", "foo", "world"}

				require.True(t, sort.Valid(whiteList))
			},
		},
		"exceptional path- field is empty": {
			validate: func(t *testing.T) {
				sort := Sort{
					Field:     "",
					Direction: "asc",
				}

				whiteList := []string{"bar", "hello", "foo", "world"}

				require.False(t, sort.Valid(whiteList))
			},
		},
		"exceptional path- field not whitelisted": {
			validate: func(t *testing.T) {
				sort := Sort{
					Field:     "foo",
					Direction: "asc",
				}

				whiteList := []string{"bar", "hello", "world"}

				require.False(t, sort.Valid(whiteList))
			},
		},
		"exceptional path- invalid direction": {
			validate: func(t *testing.T) {
				sort := Sort{
					Field:     "foo",
					Direction: "over-hill-over-dale",
				}

				whiteList := []string{"bar", "hello", "foo", "world"}

				require.False(t, sort.Valid(whiteList))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}

func TestUnit_ParsePagingParams(t *testing.T) {
	tests := map[string]struct {
		request  func(t *testing.T) *http.Request
		validate func(t *testing.T, req *http.Request)
	}{
		"base path- no defaults": {
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?offset=1&page_size=100&sort=foo:asc", nil)
				require.NoError(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				pagingParams, err := ParsePagingParams(req, PagingParams{}, []string{"foo", "bar"})
				require.NoError(t, err)

				sortFields := []Sort{
					{
						Field:     "foo",
						Direction: "asc",
					},
				}

				require.Equal(t, int32(1), *pagingParams.Offset)
				require.Equal(t, int32(100), *pagingParams.PageSize)
				require.ElementsMatch(t, sortFields, pagingParams.SortFields)
			},
		},
		"base path- all defaults set, no parameters in request": {
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com", nil)
				require.NoError(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				offset := int32(1)
				pageSize := int32(100)
				defaults := PagingParams{
					Offset:   &offset,
					PageSize: &pageSize,
					SortFields: []Sort{
						{
							Field:     "bar",
							Direction: "desc",
						},
					},
				}

				pagingParams, err := ParsePagingParams(req, defaults, []string{"foo", "bar"})
				require.NoError(t, err)

				sortFields := []Sort{
					{
						Field:     "bar",
						Direction: "desc",
					},
				}

				require.Equal(t, int32(1), *pagingParams.Offset)
				require.Equal(t, int32(100), *pagingParams.PageSize)
				require.ElementsMatch(t, sortFields, pagingParams.SortFields)
			},
		},
		"base path- default sorts provided, no sorts in request": {
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?offset=1&page_size=100", nil)
				require.NoError(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				defaults := PagingParams{
					SortFields: []Sort{
						{
							Field:     "bar",
							Direction: "desc",
						},
					},
				}

				pagingParams, err := ParsePagingParams(req, defaults, []string{"foo", "bar"})
				require.NoError(t, err)

				sortFields := []Sort{
					{
						Field:     "bar",
						Direction: "desc",
					},
				}

				require.Equal(t, int32(1), *pagingParams.Offset)
				require.Equal(t, int32(100), *pagingParams.PageSize)
				require.ElementsMatch(t, sortFields, pagingParams.SortFields)
			},
		},
		"base path- no defaults, multiple sorts given": {
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?offset=1&page_size=100&sort=foo:asc,bar:desc", nil)
				require.NoError(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				pagingParams, err := ParsePagingParams(req, PagingParams{}, []string{"foo", "bar"})
				require.NoError(t, err)

				sortFields := []Sort{
					{
						Field:     "foo",
						Direction: "asc",
					},
					{
						Field:     "bar",
						Direction: "desc",
					},
				}

				require.Equal(t, int32(1), *pagingParams.Offset)
				require.Equal(t, int32(100), *pagingParams.PageSize)
				require.ElementsMatch(t, sortFields, pagingParams.SortFields)
			},
		},
		"exceptional path- offset not an integer": {
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?offset=foo&page_size=100&sort_asc=true&sort_field=foo", nil)
				require.NoError(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				_, err := ParsePagingParams(req, PagingParams{}, []string{"foo"})
				require.Error(t, err)
			},
		},
		"exceptional path- page size not an integer": {
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?offset=100&page_size=foo&sort_asc=true&sort_field=foo", nil)
				require.NoError(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				_, err := ParsePagingParams(req, PagingParams{}, []string{"foo"})
				require.Error(t, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t, tc.request(t))
		})
	}
}
