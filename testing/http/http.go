package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

type RequestSender interface {
	Test(req *http.Request, timeout ...time.Duration) (*http.Response, error)
}

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

func Get[TSender RequestSender](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) *http.Response {
	req := MakeRequest[any](t, http.MethodGet, uri, modifiers...)

	res, err := app.Test(req)
	if err != nil {
		t.Log("Cannot get response")
		t.FailNow()
	}

	return res
}

func Post[TSender RequestSender, TBody any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) *http.Response {
	req := MakeRequest[TBody](t, http.MethodPost, uri, modifiers...)

	res, err := app.Test(req)
	if err != nil {
		t.Log("Cannot get response")
		t.FailNow()
	}

	return res
}

func Put[TSender RequestSender, TBody any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) *http.Response {
	req := MakeRequest[TBody](t, http.MethodPut, uri, modifiers...)

	res, err := app.Test(req)
	if err != nil {
		t.Log("Cannot get response")
		t.FailNow()
	}

	return res
}

func Patch[TSender RequestSender, TBody any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) *http.Response {
	req := MakeRequest[TBody](t, http.MethodPatch, uri, modifiers...)

	res, err := app.Test(req)
	if err != nil {
		t.Log("Cannot get response")
		t.FailNow()
	}

	return res
}

func Delete[TSender RequestSender](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) *http.Response {
	req := MakeRequest[any](t, http.MethodDelete, uri, modifiers...)

	res, err := app.Test(req)
	if err != nil {
		t.Log("Cannot get response")
		t.FailNow()
	}

	return res
}
