package http

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestWithQuerySuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)
	qs := map[string]string{
		"id": "12345",
	}

	httpReq, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/click", nil)

	type Test struct {
		Id int64
	}

	// Act
	req := WithQuery[Test](t, qs)

	// Assert
	assert.Equal("12345", req(httpReq).URL.Query().Get("id"))
}

func TestName(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)
	headers := http.Header{}
	headers.Set("Content-Type", "")
	req, _ := http.NewRequest(http.MethodGet, "", nil)

	// Act
	r := WithHeaders(headers)

	// Assert
	assert.Equal(r(req).Header.Get("Content-Type"), "application/json")
}

func TestWithCookie(t *testing.T) {
	// Assert
	t.Parallel()
	assert := require.New(t)
	cookies := []*http.Cookie{
		{Name: "test"},
	}
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	// Act
	res := WithCookies(cookies)
	// Assert
	a, _ := res(req).Cookie("test")
	assert.Equal(a.Name, "test")
}
