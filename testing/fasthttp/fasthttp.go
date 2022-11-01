package fasthttp

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

// Type FastHttpSender
// Sturcture contains in memory server and client for testing purposes
type FastHttpSender[T any] struct {
	app            *fasthttp.Server
	testing        testing.TB
	followRedirect bool
	ln             *fasthttputil.InmemoryListener
}

// Instantiate new Fast HTTP Client for testing purposes
func New[T any](t testing.TB, app *fasthttp.Server, followRedirects ...bool) *FastHttpSender[T] {
	ln := fasthttputil.NewInmemoryListener()

	var followRedirect bool

	if len(followRedirects) > 0 {
		followRedirect = followRedirects[0]
	}

	go func() {
		if err := app.Serve(ln); err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}()

	sender := &FastHttpSender[T]{
		ln:             ln,
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

// Sends a HTTP request for testing purposes
func (s *FastHttpSender[T]) Test(req *http.Request, timeout ...time.Duration) (*http.Response, error) {
	var tmout time.Duration

	if len(timeout) > 0 {
		tmout = timeout[0]
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return s.ln.Dial()
			},
		},
		Timeout: tmout,
	}

	if !s.followRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	res, err := client.Do(req)

	if err != nil {
		s.testing.Log(err)
		s.testing.FailNow()
	}

	return res, nil
}

func (s FastHttpSender[T]) Close() error {
	return s.app.Shutdown()
}
