package fasthttp

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	nano_http "github.com/nano-interactive/go-utils/v2/testing/http"
	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/valyala/fasthttp"
)

type Sender[T any] struct {
	dialer         func() (net.Conn, error)
	testing        testing.TB
	followRedirect bool
}

func NewWithServer[T any](t testing.TB, s *fasthttp.Server, followRedirects ...bool) *Sender[T] {
	t.Helper()

	ln := fasthttputil.NewInmemoryListener()

	go func() {
		if err := s.Serve(ln); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}()

	t.Cleanup(func() {
		if err := ln.Close(); err != nil {
			t.Errorf("Failed to close FastHTTP Server")
		}
	})

	return New[T](t, ln.Dial, followRedirects...)
}

func New[T any](t testing.TB, dialer func() (net.Conn, error), followRedirects ...bool) *Sender[T] {
	t.Helper()
	var followRedirect bool

	if len(followRedirects) > 0 {
		followRedirect = followRedirects[0]
	}

	sender := &Sender[T]{
		dialer:         dialer,
		testing:        t,
		followRedirect: followRedirect,
	}

	return sender
}

// Sends a HTTP request for testing purposes
func (s *Sender[T]) Test(req *http.Request, timeout ...time.Duration) nano_http.ExtendedResponse[T] {
	var tmout time.Duration

	if len(timeout) > 0 {
		tmout = timeout[0]
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return s.dialer()
			},
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return s.dialer()
			},
			DisableKeepAlives:  true,
			DisableCompression: true,
			MaxIdleConns:       1,
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
		s.testing.Errorf("failed to send request: %v", err)
		s.testing.FailNow()
	}

	return nano_http.ExtendedResponse[T]{Response: res}
}

func QueryArgsFromUrl(url string) *fasthttp.Args {
	args := &fasthttp.Args{}

	args.Parse(url)

	return args
}

func CallHandler(t testing.TB, url string, h fasthttp.RequestHandler) *fasthttp.RequestCtx {
	t.Helper()

	u := fasthttp.URI{}
	u.Update(url)

	fastHttpCtx := &fasthttp.RequestCtx{
		Request: fasthttp.Request{
			Header:        fasthttp.RequestHeader{},
			UseHostHeader: false,
		},
		Response: fasthttp.Response{},
	}
	fastHttpCtx.Request.SetURI(&u)

	h(fastHttpCtx)

	return fastHttpCtx
}
