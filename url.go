package utils

import (
	"bytes"
	"errors"
	"net/url"
	"strings"
)

var (
	GooglePlay  = []byte("play.google")
	ITunesApple = []byte("itunes.apple")

	HTTPS         = []byte("https://")
	HTTP          = []byte("http://")
	HTTPSProtocol = []byte("https:")
	HTTPProtocol  = []byte("http:")

	ErrInvalidUrl = errors.New("invalid url")
)

func TrimUrlForScylla(fullUrl string) (scyllaUrl string, hostName string, err error) {
	var data strings.Builder
	trimmedUrl := strings.TrimSpace(fullUrl)

	urlObject, err := url.Parse(trimmedUrl)
	if err != nil {
		//remove everything behind ? or #
		if index := strings.Index(trimmedUrl, "?"); index > 0 {
			trimmedUrl = trimmedUrl[0:index]
		}
		if index := strings.Index(trimmedUrl, "#"); index > 0 {
			trimmedUrl = trimmedUrl[0:index]
		}
		trimmedUrl = strings.ToValidUTF8(trimmedUrl, "")

		urlObject, err = url.Parse(trimmedUrl)
		if err != nil {
			return "", "", err
		}
	}

	data.Grow(len(urlObject.Host) + len(urlObject.Path) + len("https://") + 1)
	data.WriteString("https://")
	data.WriteString(urlObject.Host)
	data.WriteString(urlObject.Path)
	scyllaUrl = data.String()

	if scyllaUrl[len(scyllaUrl)-1:] != "/" {
		data.WriteString("/")
	}
	scyllaUrl = data.String()
	scyllaUrl = strings.ToValidUTF8(scyllaUrl, "")

	return scyllaUrl, urlObject.Host, nil
}

func GetDomainFromUrl(fullUrl string) (string, error) {
	fullUrlObject, err := url.Parse(fullUrl)
	if err != nil {
		return "", err
	}
	return fullUrlObject.Hostname(), nil
}

// Truncates a given url
func TruncateUrl(value []byte) []byte {
	if bytes.Contains(value, GooglePlay) || bytes.Contains(value, ITunesApple) {
		return value
	}

	n := bytes.IndexAny(value, "?#")
	if n != -1 {
		value = value[0:n]
	}

	if value[len(value)-1] != byte('/') {
		value = append(value, byte('/'))
	}

	if !bytes.Equal(value[0:4], []byte("http")) {
		value = append(HTTPS, value...)
	}

	if bytes.Equal(value[0:5], HTTPProtocol) {
		value = append(HTTPSProtocol, value[5:]...)
	}

	return value
}
