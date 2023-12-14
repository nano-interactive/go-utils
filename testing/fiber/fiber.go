package fiber

import (
	"net/http"
	"testing"
	"time"

	nano_http "github.com/nano-interactive/go-utils/v2/testing/http"
	"github.com/valyala/fasthttp"
)

// type GoFiberSender
// Sturcture contains in memory server and client for testing purposes
type GoFiberSender[T any] struct {
	app            *fiber.App
	testing        testing.TB
	followRedirect bool
}

// Instantiate New fiber client for testing purposes
func New[T any](t testing.TB, app *fiber.App, followRedirects ...bool) *GoFiberSender[T] {
	t.Helper()

	var followRedirect bool

	if len(followRedirects) > 0 {
		followRedirect = followRedirects[0]
	}

	sender := &GoFiberSender[T]{
		app:            app,
		testing:        t,
		followRedirect: followRedirect,
	}

	t.Cleanup(func() {
		if err := app.Shutdown(); err != nil {
			t.Errorf("failed to close sender and server: %v", err)
		}
	})

	return sender
}

// Sends a new Fiber request for testing purposes
func (s *GoFiberSender[T]) Test(req *http.Request, timeout ...time.Duration) nano_http.ExtendedResponse[T] {
	s.testing.Helper()

	t := -1

	if len(timeout) > 0 {
		t = int(timeout[0].Seconds())
	}

	res, err := s.app.Test(req, t)
	if err != nil {
		s.testing.Errorf("Failed to DO request: %v", err)
		return nano_http.ExtendedResponse[T]{}
	}

	return nano_http.ExtendedResponse[T]{Response: res}
}

func CallHandler(tb testing.TB, h fiber.Handler, fn ...func(*fasthttp.Request)) (*fiber.Ctx, error) {
	tb.Helper()

	app := fiber.New()

	fastHTTPCtx := &fasthttp.RequestCtx{
		Request: fasthttp.Request{
			Header:        fasthttp.RequestHeader{},
			UseHostHeader: false,
		},
		Response: fasthttp.Response{},
	}

	if len(fn) > 0 {
		for _, f := range fn {
			f(&fastHTTPCtx.Request)
		}
	}

	ctx := app.AcquireCtx(fastHTTPCtx)

	tb.Cleanup(func() {
		app.ReleaseCtx(ctx)
	})

	return ctx, h(ctx)
}
