package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"testing"
	"time"
)

type RequestSender[T any] interface {
	Test(req *http.Request, timeout ...time.Duration) ExtendedResponse[T]
}

// Returns io.Reader from Body of test request
func getBody[T any](t testing.TB, headers http.Header, body T) io.Reader {
	switch headers.Get("Content-Type") {
	case fiber.MIMEApplicationXML, fiber.MIMEApplicationXMLCharsetUTF8:
		bs, err := xml.Marshal(body)
		if err != nil {
			t.Errorf("Error while creating XML from Body: %v", err)
			t.FailNow()
		}

		return bytes.NewReader(bs)
	case fiber.MIMEApplicationJSON, fiber.MIMEApplicationJSONCharsetUTF8:
		bs, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Error while sending request: %v", err)
			t.FailNow()
		}

		return bytes.NewReader(bs)
	default:
		return nil
	}
}

func Get[TSender RequestSender[any]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[any] {
	return app.Test(MakeRequest(t, http.MethodGet, uri, modifiers...))
}

func GetWithResponse[TSender RequestSender[TRes], TRes any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[TRes] {
	return app.Test(MakeRequest(t, http.MethodGet, uri, modifiers...))
}

func Post[TSender RequestSender[any]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[any] {
	return app.Test(MakeRequest(t, http.MethodPost, uri, modifiers...))
}

func PostWithResponse[TSender RequestSender[TRes], TRes any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[TRes] {
	return app.Test(MakeRequest(t, http.MethodPost, uri, modifiers...))
}

func Put[TSender RequestSender[any]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[any] {
	return app.Test(MakeRequest(t, http.MethodPut, uri, modifiers...))
}

func PutWithResponse[TSender RequestSender[TRes], TRes any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[TRes] {
	return app.Test(MakeRequest(t, http.MethodPut, uri, modifiers...))
}

func Patch[TSender RequestSender[any]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[any] {
	return app.Test(MakeRequest(t, http.MethodPatch, uri, modifiers...))
}

func PatchWithResponse[TSender RequestSender[TRes], TRes any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[TRes] {
	return app.Test(MakeRequest(t, http.MethodPatch, uri, modifiers...))
}

func Delete[TSender RequestSender[any]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[any] {
	return app.Test(MakeRequest(t, http.MethodDelete, uri, modifiers...))
}

func DeleteWithResponbe[TSender RequestSender[TRes], TRes any](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[TRes] {
	return app.Test(MakeRequest(t, http.MethodDelete, uri, modifiers...))
}
