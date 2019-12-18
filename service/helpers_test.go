package service

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnit_Int32PointerFromQueryParam(t *testing.T) {

	type testCase struct {
		name     string
		request  func(t *testing.T) *http.Request
		validate func(t *testing.T, req *http.Request)
	}

	tests := []testCase{
		{
			name: "param exists and has a castable value",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=1", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := Int32PointerFromQueryParam(req, "foo")
				require.Nil(t, err)

				require.Equal(t, int32(1), *result)
			},
		},
		{
			name: "param exists and has empty value - returns nil int32",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := Int32PointerFromQueryParam(req, "foo")
				require.Nil(t, err)
				require.Nil(t, result)
			},
		},
		{
			name: "param exists and has a non-castable value - returns error",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=bar", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				_, err := Int32PointerFromQueryParam(req, "foo")
				require.NotNil(t, err)
			},
		},
		{
			name: "param does not exist - returns nil int32",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?nope=value", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := Int32PointerFromQueryParam(req, "foo")
				require.Nil(t, err)
				require.Nil(t, result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.validate(t, tc.request(t))
		})
	}
}

func TestUnit_Int64ArrayFromQueryParam(t *testing.T) {

	type testCase struct {
		name     string
		request  func(t *testing.T) *http.Request
		validate func(t *testing.T, req *http.Request)
	}

	tests := []testCase{
		{
			name: "param exists and has a castable value",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=1,2,3", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := Int64ArrayFromQueryParam(req, "foo")
				require.Nil(t, err)

				require.Equal(t, []int64{1, 2, 3}, result)
			},
		},
		{
			name: "param exists and has empty value - returns nil",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := Int64ArrayFromQueryParam(req, "foo")
				require.Nil(t, err)
				require.Nil(t, result)
			},
		},
		{
			name: "param exists and has a non-castable value - returns error",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=bar", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				_, err := Int64ArrayFromQueryParam(req, "foo")
				require.NotNil(t, err)
			},
		},
		{
			name: "param does not exist - returns nil",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?nope=value", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := Int64ArrayFromQueryParam(req, "foo")
				require.Nil(t, err)
				require.Nil(t, result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.validate(t, tc.request(t))
		})
	}
}

func TestUnit_TimestampFromQueryParam(t *testing.T) {

	type testCase struct {
		name     string
		request  func(t *testing.T) *http.Request
		validate func(t *testing.T, req *http.Request)
	}

	tests := []testCase{
		{
			name: "param exists and has a castable value",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=2006-01-02T15:04:05Z", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := TimestampFromQueryParam(req, "foo")
				require.Nil(t, err)

				expected, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
				require.Equal(t, expected, *result)
			},
		},
		{
			name: "param exists and has empty value - returns nil",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := TimestampFromQueryParam(req, "foo")
				require.Nil(t, err)
				require.Nil(t, result)
			},
		},
		{
			name: "param exists and has a non-castable value - returns error",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=Mon, 02 Jan 2006 15:04:05 MST", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				_, err := TimestampFromQueryParam(req, "foo")
				require.NotNil(t, err)
			},
		},
		{
			name: "param does not exist - returns nil",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?nope=value", nil)
				require.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := TimestampFromQueryParam(req, "foo")
				require.Nil(t, err)
				require.Nil(t, result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.validate(t, tc.request(t))
		})
	}
}
