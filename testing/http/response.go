package http

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type ExtendedResponse[T any] struct {
	*http.Response
}

func (e ExtendedResponse[T]) AssertStatusCode(t testing.TB, status int) {
	if e.StatusCode != status {
		t.Errorf("Expected HTTP Status Code %d(%s), Actual %d(%s)", status, http.StatusText(status), e.StatusCode, e.Status)
		t.Fail()
	}
}

func (e ExtendedResponse[T]) AssertHeader(t testing.TB, header, value string) {
	val := e.Header.Get(header)

	if val == "" && value != "" {
		t.Errorf("Expected HTTP Header to be %s, but it's empty", value)
		t.Fail()
	}

	if val != value {
		t.Errorf("Expected HTTP Header to be %s, but it's %s", value, val)
		t.Fail()
	}
}

func (e ExtendedResponse[T]) Data(tb testing.TB) T {
	tb.Helper()

	tb.Cleanup(func() {
		if err := e.Response.Body.Close(); err != nil {
			tb.Errorf("Failed to close Response BODY: %v", err)
		}
	})

	var data T

	switch e.Response.Header.Get(fiber.HeaderContentType) {
	case fiber.MIMEApplicationXML, fiber.MIMEApplicationXMLCharsetUTF8:
		if err := xml.NewDecoder(e.Response.Body).Decode(&data); err != nil {
			tb.Errorf("Invalid XML Response: %v", err)
		}
	case "":
		tb.Log("Content Type header is not set => using JSON as Default")
		fallthrough
	case fiber.MIMEApplicationJSON, fiber.MIMEApplicationJSONCharsetUTF8:
		fallthrough
	default:
		if err := json.NewDecoder(e.Response.Body).Decode(&data); err != nil {
			tb.Errorf("Invalid JSON Response: %v", err)
		}
	}

	return data
}
