package httpu_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/nano-interactive/go-utils/v2/httpu"
	"github.com/stretchr/testify/assert"
)

func TestCurlFromRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *http.Request
		expected string
		err      error
	}{
		{
			name: "Valid GET request",
			request: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/"},
				Header: http.Header{
					"User-Agent": []string{"test-agent"},
				},
			},
			expected: "curl -s -v -X GET -H 'User-Agent: test-agent' --compressed 'http://example.com/'",
			err:      nil,
		},
		{
			name: "Valid POST request with body",
			request: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/submit"},
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(bytes.NewBufferString(`{"key":"value"}`)),
			},
			expected: "curl -s -v -X POST -H 'Content-Type: application/json' -d '{\"key\":\"value\"}' --compressed 'http://example.com/submit'",
			err:      nil,
		},
		{
			name: "Request with no URL",
			request: &http.Request{
				Method: http.MethodGet,
				URL:    nil,
			},
			expected: "",
			err:      httpu.ErrNoRequestURL,
		},
		{
			name: "Request with multiple headers",
			request: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/"},
				Header: http.Header{
					"User-Agent": []string{"test-agent"},
					"Accept":     []string{"application/json", "text/html"},
				},
			},
			expected: "curl -s -v -X GET -H 'User-Agent: test-agent' -H 'Accept: application/json' -H 'Accept: text/html' --compressed 'http://example.com/'",
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := httpu.CurlFromRequest(tt.request)

			if tt.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
