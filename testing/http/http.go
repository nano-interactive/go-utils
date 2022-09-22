package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

// Type RequestSender
// Sturcture contains in memory HTTP Server and Client for testing purposes
type RequestSender interface {
	Test(req *http.Request, timeout ...time.Duration) (*http.Response, error)
}

// Returns io.Reader from Body of test request
func getBody[T any](t testing.TB, headers http.Header, body T) io.Reader {
	switch headers.Get("Content-Type") {
	case "application/json":
		jsonStr, err := json.Marshal(body)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		return bytes.NewReader(jsonStr)
	default:
		return nil
	}
}

// Returns pointer to http.Response from GET request for testing purposes
func Get[TSender RequestSender](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) *http.Response {
	req := MakeRequest[any](t, http.MethodGet, uri, modifiers...)

	res, err := app.Test(req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	return res
}

// Returns pointer to http.Response from POST request for testing purposes
func Post[TSender RequestSender, TBody any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) *http.Response {
	req := MakeRequest[TBody](t, http.MethodPost, uri, modifiers...)

	res, err := app.Test(req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	return res
}

// Returns pointer to http.Response from PUT request for testing purposes
func Put[TSender RequestSender, TBody any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) *http.Response {
	req := MakeRequest[TBody](t, http.MethodPut, uri, modifiers...)

	res, err := app.Test(req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	return res
}

// Returns pointer to http.Response from PATCH request for testing purposes
func Patch[TSender RequestSender, TBody any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) *http.Response {
	req := MakeRequest[TBody](t, http.MethodPatch, uri, modifiers...)

	res, err := app.Test(req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	return res
}

// Returns pointer to http.Response from DELETE request for testing purposes
func Delete[TSender RequestSender](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) *http.Response {
	req := MakeRequest[any](t, http.MethodDelete, uri, modifiers...)

	res, err := app.Test(req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	return res
}
