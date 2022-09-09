package fiber

import (
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GoFiberSender[T any] struct {
	app            *fiber.App
	testing        testing.TB
	followRedirect bool
}

func New[T any](t testing.TB, app *fiber.App, followRedirect bool) *GoFiberSender[T] {
	return &GoFiberSender[T]{
		app:            app,
		testing:        t,
		followRedirect: followRedirect,
	}
}

func (s *GoFiberSender[T]) Test(req *http.Request, timeout ...time.Duration) (*http.Response, error) {
	t := -1

	if len(timeout) > 0 {
		t = int(timeout[0].Seconds())
	}

	return s.app.Test(req, t)
}
