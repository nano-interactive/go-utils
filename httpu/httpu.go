package httpu

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/nano-interactive/go-utils/v2/constants"
)

var (
	ErrNoRequestURL = errors.New("request url is empty")
)

func CurlFromRequest(request *http.Request) (string, error) {
	if request.URL == nil {
		return "", ErrNoRequestURL
	}

	var b strings.Builder
	if request.Body == nil {
		b.Grow(1 * constants.KiB)
	} else {
		// It's better to overestimate
		// the amount of required memory
		// than to not allocate enough.
		b.Grow(1 * constants.MiB)
	}

	b.WriteString("curl -s -v -X ")
	b.WriteString(request.Method)
	b.WriteRune(' ')

	for headerName, headerValues := range request.Header {
		for _, headerValue := range headerValues {
			b.WriteString("-H '")
			b.WriteString(headerName)
			b.WriteString(": ")
			b.WriteString(headerValue)
			b.WriteString("' ")
		}
	}

	if request.Body != nil {
		var buff bytes.Buffer
		_, err := buff.ReadFrom(request.Body)
		if err != nil {
			return "", err
		}

		b.WriteString("-d '")
		b.WriteString(buff.String())
		b.WriteString("' ")

		// reset the body after reading
		request.Body = io.NopCloser(bytes.NewBuffer(buff.Bytes()))
	}

	b.WriteString("--compressed '")
	b.WriteString(request.URL.String())
	b.WriteString("'")

	return b.String(), nil
}
