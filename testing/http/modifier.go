package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type RequestModifier func(*http.Request) *http.Request

func WithHeaders(t testing.TB, headers http.Header) RequestModifier {
	return func(req *http.Request) *http.Request {
		if headers.Get("Content-Type") == "" {
			headers.Set("Content-Type", "application/json")
		}

		if headers.Get("Accept") == "" {
			headers.Set("Accept", "application/json")
		}

		if headers.Get("User-Agent") == "" {
			headers.Set("User-Agent", "TestHTTPUserAgent")
		}

		req.Header = headers

		return req
	}
}

// WithQuery Provides a way to put query string params.
//
//	WithQuery(map[string]string{
//		"id": "mongoid",
//	})
//
// Returns URL?id=mongoid
func WithQuery[T any](t testing.TB, queryMap map[string]string) RequestModifier {
	return func(req *http.Request) *http.Request {
		newReq := httptest.NewRequest(req.Method, req.URL.Path, nil)
		newReq.Header = req.Header
		query := newReq.URL.Query()

		for key, value := range queryMap {
			query.Add(key, value)
		}

		newReq.URL.RawQuery = query.Encode()

		return newReq
	}
}

func WithBody[T any](t testing.TB, body T) RequestModifier {
	return func(req *http.Request) *http.Request {
		newReq := httptest.NewRequest(req.Method, req.URL.Path, getBody[T](t, req.Header, body))
		newReq.Header = req.Header
		newReq.URL.RawQuery = req.URL.Query().Encode()

		for _, cookie := range req.Cookies() {
			newReq.AddCookie(cookie)
		}

		return newReq
	}
}

func WithCookies(t testing.TB, cookies []*http.Cookie) RequestModifier {
	return func(req *http.Request) *http.Request {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		return req
	}
}

func MakeRequest[T any](t testing.TB, method, uri string, modifiers ...RequestModifier) *http.Request {
	var defaults []func(*http.Request) *http.Request

	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		defaults = []func(*http.Request) *http.Request{
			WithHeaders(t, http.Header{}),
			WithBody[any](t, nil),
		}
	default:
		defaults = []func(*http.Request) *http.Request{
			WithHeaders(t, http.Header{}),
		}
	}

	req := httptest.NewRequest(method, uri, nil)

	for _, modifier := range defaults {
		req = modifier(req)
	}

	for _, modifier := range modifiers {
		req = modifier(req)
	}

	return req
}
