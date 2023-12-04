package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RequestSender[T any] interface {
	Test(req *http.Request, timeout ...time.Duration) ExtendedResponse[T]
}

// Returns io.Reader from Body of test request
func getBody[T any](t testing.TB, headers http.Header, body T) io.Reader {
	contentType := headers.Get(fiber.HeaderContentType)
	switch contentType {
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
		var body1 any = body
		if body1 == nil {
			return nil
		}
		return body1.(io.Reader)
	}
}

func Get[TSender RequestSender[any]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[any] {
	return app.Test(MakeRequest(t, http.MethodGet, uri, modifiers...))
}

func GetWithResponse[TRes any, TSender RequestSender[TRes]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[TRes] {
	return app.Test(MakeRequest(t, http.MethodGet, uri, modifiers...))
}

func Post[TSender RequestSender[any]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[any] {
	return app.Test(MakeRequest(t, http.MethodPost, uri, modifiers...))
}

func PostWithResponse[TRes any, TSender RequestSender[TRes]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[TRes] {
	return app.Test(MakeRequest(t, http.MethodPost, uri, modifiers...))
}

func Put[TSender RequestSender[any]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[any] {
	return app.Test(MakeRequest(t, http.MethodPut, uri, modifiers...))
}

func PutWithResponse[TRes any, TSender RequestSender[TRes]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[TRes] {
	return app.Test(MakeRequest(t, http.MethodPut, uri, modifiers...))
}

func Patch[TSender RequestSender[any]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[any] {
	return app.Test(MakeRequest(t, http.MethodPatch, uri, modifiers...))
}

func PatchWithResponse[TRes any, TSender RequestSender[TRes]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[TRes] {
	return app.Test(MakeRequest(t, http.MethodPatch, uri, modifiers...))
}

func Delete[TSender RequestSender[any]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[any] {
	return app.Test(MakeRequest(t, http.MethodDelete, uri, modifiers...))
}

func DeleteWithResponbe[TRes any, TSender RequestSender[TRes]](t testing.TB, app TSender, uri string, modifiers ...RequestModifier) ExtendedResponse[TRes] {
	return app.Test(MakeRequest(t, http.MethodDelete, uri, modifiers...))
}
