package fiber

import (
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// type GoFiberSender
// Sturcture contains in memory server and client for testing purposes
type GoFiberSender struct {
	app            *fiber.App
	testing        testing.TB
	followRedirect bool
}

// Instantiate New fiber client for testing purposes
func New(t testing.TB, app *fiber.App, followRedirects ...bool) *GoFiberSender {
	var followRedirect bool

	if len(followRedirects) > 0 {
		followRedirect = followRedirects[0]
	}

	sender := &GoFiberSender{
		app:            app,
		testing:        t,
		followRedirect: followRedirect,
	}

	t.Cleanup(func() {
		err := sender.Close()
		if err != nil {
			t.Fatalf("failed to close sender and server: %v", err)
		}
	})

	return sender
}

// Sends a new Fiber request for testing purposes
func (s *GoFiberSender) Test(req *http.Request, timeout ...time.Duration) (*http.Response, error) {
	t := -1

	if len(timeout) > 0 {
		t = int(timeout[0].Seconds())
	}

	return s.app.Test(req, t)
}

func (s *GoFiberSender) Close() error {
	return s.app.Shutdown()
}
