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
}

// Instantiate new Fast HTTP Client for testing purposes
func New[T any](t testing.TB, app *fasthttp.Server, followRedirect bool) *FastHttpSender[T] {
	return &FastHttpSender[T]{
		app:            app,
		testing:        t,
		followRedirect: followRedirect,
	}
}

// Sends a HTTP request for testing purposes
func (s *FastHttpSender[T]) Test(req *http.Request, timeout ...time.Duration) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()

	go func() {
		if err := s.app.Serve(ln); err != nil {
			s.testing.Log(err)
			s.testing.FailNow()
		}
	}()

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	if s.followRedirect {
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
