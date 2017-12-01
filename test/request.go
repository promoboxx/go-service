package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewJSONRequest will return a request with body marshalled into json
func NewJSONRequest(ctx context.Context, t *testing.T, method, target string, body interface{}) *http.Request {
	var by []byte
	var ok bool
	var err error
	if by, ok = body.([]byte); !ok {
		by, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Could not json marshal body: %v", err)
		}
	}

	return httptest.NewRequest(method, target, bytes.NewBuffer(by)).WithContext(ctx)
}
