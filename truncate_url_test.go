package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var data1 = []string{
	"https://www.example.com/",                               // valid
	"https://www.example.com",                                // missing the last slash
	"http://www.example.com",                                 // missing https
	"www.example.com/",                                       // missing protocol
	"www.example.com",                                        // missing protocol and slash
	"www.example.com?one=two",                                // query string
	"www.example.com#one=two",                                // query string
	"www.example.com?one=two&two=http",                       // query string no protocol
	"https://www.example.com?one=two&two=http",               // query string two query params, no /
	"https://www.example.com/?one=two&two=http",              // query string two query params, with /
	"http://www.example.com?one=two&two=https://yup.com?he",  // query string without / and with weird query param
	"http://www.example.com/?one=two&two=https://yup.com?he", // query string with / and with weird query param
}

var data2 = []string{
	"https://www.example.com/data/",                               // valid
	"https://www.example.com/data",                                // missing the last slash
	"http://www.example.com/data",                                 // missing https
	"www.example.com/data/",                                       // missing protocl
	"www.example.com/data",                                        // missing protocol and slash
	"www.example.com/data?one=two",                                // query string
	"www.example.com/data#one=two",                                // query string with #
	"www.example.com/data?one=two&two=http",                       // query string no protocol
	"https://www.example.com/data?one=two&two=http",               // query string two query params, no /
	"https://www.example.com/data/?one=two&two=http",              // query string two query params, with /
	"http://www.example.com/data?one=two&two=https://yup.com?he",  // query string without / and with wierd query param
	"http://www.example.com/data/?one=two&two=https://yup.com?he", // query string with / and with wierd query param
}

func TestTruncateUrlSimple(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	for _, url := range data1 {
		result := TruncateUrl([]byte(url))
		assert.Equal([]byte("https://www.example.com/"), result, "URL not matches")
	}
}

func TestTruncateUrlDeepLink(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	for _, url := range data2 {
		result := TruncateUrl([]byte(url))
		assert.Equal([]byte("https://www.example.com/data/"), result, "URL not matches")
	}
}

func BenchmarkTruncateUrl(b *testing.B) {
	for _, url := range data2 {
		u := []byte(url)
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			TruncateUrl(u)
		}
	}
}
