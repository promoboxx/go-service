package service

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := Int32PointerFromQueryParam(req, "foo")
				assert.Nil(t, err)

				assert.Equal(t, int32(1), *result)
			},
		},
		{
			name: "param exists and has empty value - returns nil int32",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=", nil)
				assert.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := Int32PointerFromQueryParam(req, "foo")
				assert.Nil(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "param exists and has a non-castable value - returns error",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?foo=bar", nil)
				assert.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				_, err := Int32PointerFromQueryParam(req, "foo")
				assert.NotNil(t, err)
			},
		},
		{
			name: "param does not exist - returns nil int32",
			request: func(t *testing.T) *http.Request {
				req, err := http.NewRequest("GET", "http://example.com?nope=value", nil)
				assert.Nil(t, err)
				return req
			},
			validate: func(t *testing.T, req *http.Request) {
				result, err := Int32PointerFromQueryParam(req, "foo")
				assert.Nil(t, err)
				assert.Nil(t, result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.validate(t, tc.request(t))
		})
	}
}
