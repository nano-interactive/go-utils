package fasthttp

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

// Type FastHttpSender
// Sturcture contains in memory server and client for testing purposes
type FastHttpSender struct {
	dialer         func() (net.Conn, error)
	testing        testing.TB
	followRedirect bool
}

// Instantiate new Fast HTTP Client for testing purposes
func New(t testing.TB, dialer func() (net.Conn, error), followRedirects ...bool) *FastHttpSender {
	var followRedirect bool

	if len(followRedirects) > 0 {
		followRedirect = followRedirects[0]
	}

	sender := &FastHttpSender{
		dialer:         dialer,
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
func (s *FastHttpSender) Test(req *http.Request, timeout ...time.Duration) (*http.Response, error) {
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
		s.testing.Fatalf("failed to send request: %v", err)
	}

	return res, nil
}

func (s FastHttpSender) Close() error {
	return nil
}

func QueryArgsFromUrl(url string) *fasthttp.Args {
	args := &fasthttp.Args{}

	args.Parse(url)

	return args
}

func CallHandler(t testing.TB, url string, h fasthttp.RequestHandler) *fasthttp.RequestCtx {
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
