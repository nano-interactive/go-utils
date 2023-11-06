package http

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type RequestModifier func(*http.Request) *http.Request

func WithBearerToken(t testing.TB, token string) RequestModifier {
	t.Helper()

	return func(req *http.Request) *http.Request {
		req.Header.Add(fiber.HeaderAuthorization, "Bearer "+token)

		return req
	}
}

func WithHeaders(t testing.TB, headers http.Header) RequestModifier {
	t.Helper()
	return func(req *http.Request) *http.Request {
		if headers.Get(fiber.HeaderContentType) == "" {
			headers.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		}

		if headers.Get(fiber.HeaderAccept) == "" {
			headers.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSONCharsetUTF8)
		}

		if headers.Get(fiber.HeaderUserAgent) == "" {
			headers.Set(fiber.HeaderUserAgent, "TestHTTPUserAgent")
		}

		req.Header = headers

		return req
	}
}

func WithQuery(t testing.TB, queryMap map[string]string) RequestModifier {
	t.Helper()
	return func(req *http.Request) *http.Request {
		newReq, err := http.NewRequest(req.Method, req.URL.String(), nil)
		if err != nil {
			t.Errorf("Failed to create new http.Request: %v", err)
			t.FailNow()
		}

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
	t.Helper()
	return func(req *http.Request) *http.Request {
		newReq, err := http.NewRequest(req.Method, req.URL.String(), getBody(t, req.Header, body))
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		newReq.Header = req.Header
		newReq.URL.RawQuery = req.URL.Query().Encode()

		for _, cookie := range req.Cookies() {
			newReq.AddCookie(cookie)
		}

		return newReq
	}
}

func WithCookies(t testing.TB, cookies []*http.Cookie) RequestModifier {
	t.Helper()
	return func(req *http.Request) *http.Request {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		return req
	}
}

func MakeRequest(t testing.TB, method, uri string, modifiers ...RequestModifier) *http.Request {
	t.Helper()
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

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, modifier := range defaults {
		req = modifier(req)
	}

	for _, modifier := range modifiers {
		req = modifier(req)
	}

	return req
}
