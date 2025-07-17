package utils

import (
	"bytes"
	"errors"
	net_url "net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
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
	trimmedUrl := strings.TrimSpace(fullUrl)
	if len(trimmedUrl) == 0 {
		return "", "", ErrInvalidUrl
	}
	// Remove everything behind ? or #
	if index := strings.Index(trimmedUrl, "?"); index > 0 {
		trimmedUrl = trimmedUrl[0:index]
	}
	if index := strings.Index(trimmedUrl, "#"); index > 0 {
		trimmedUrl = trimmedUrl[0:index]
	}
	trimmedUrl = strings.ToValidUTF8(trimmedUrl, "")

	// Remove trailing './' if present
	trimmedUrl = strings.TrimSuffix(trimmedUrl, ".")
	trimmedUrl = strings.TrimSuffix(trimmedUrl, "./")

	if trimmedUrl[len(trimmedUrl)-1:] != "/" {
		trimmedUrl = trimmedUrl + "/"
	}

	if strings.Index(trimmedUrl, "http") != 0 {
		trimmedUrl = "https://" + trimmedUrl
	} else {
		trimmedUrl = strings.Replace(trimmedUrl, "http://", "https://", 1)
	}

	if len(trimmedUrl) < 11 {
		return "", "", ErrInvalidUrl
	}
	host := trimmedUrl[8:]
	if index := strings.Index(host, "/"); index > 0 {
		host = host[0:index]
	}

	return trimmedUrl, host, nil
}

func TrimUrlForScyllaOld(fullUrl string) (scyllaUrl string, hostName string, err error) {
	var data strings.Builder
	trimmedUrl := strings.TrimSpace(fullUrl)

	urlObject, err := net_url.Parse(trimmedUrl)
	if err != nil {
		// remove everything behind ? or #
		if index := strings.Index(trimmedUrl, "?"); index > 0 {
			trimmedUrl = trimmedUrl[0:index]
		}
		if index := strings.Index(trimmedUrl, "#"); index > 0 {
			trimmedUrl = trimmedUrl[0:index]
		}
		trimmedUrl = strings.ToValidUTF8(trimmedUrl, "")

		urlObject, err = net_url.Parse(trimmedUrl)
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
	fullUrlObject, err := net_url.Parse(fullUrl)
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

func CleanDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)

	// Remove scheme if present
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	// Remove "www." prefix and trailing slash
	domain = strings.TrimPrefix(domain, "www.")
	domain = strings.TrimSuffix(domain, "/")

	return domain
}

func ExtractTldPlusOne(url string) (string, error) {
	if url == "" {
		return "", errors.New("empty url")
	}

	var host string
	url = strings.ReplaceAll(url, "http://", "https://")

	if !strings.Contains(url, "://") {
		url = "https://" + url
	}

	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	host = url
	parsed, err := net_url.Parse(url)
	if err == nil {
		host = parsed.Hostname()
	}

	eTLDPlusOne, errEtldError := publicsuffix.EffectiveTLDPlusOne(host)

	eTLDPlusOne = strings.ToLower(eTLDPlusOne)
	if strings.HasPrefix(eTLDPlusOne, "www.") {
		eTLDPlusOne = strings.TrimPrefix(eTLDPlusOne, "www.")
	}
	if strings.HasPrefix(eTLDPlusOne, "https://") {
		eTLDPlusOne = strings.TrimPrefix(eTLDPlusOne, "https://")
	}

	return eTLDPlusOne, errEtldError
}
