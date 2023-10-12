package utils

import (
	"net/url"
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
	"www.example.com/data/",                                       // missing protocol
	"www.example.com/data",                                        // missing protocol and slash
	"www.example.com/data?one=two",                                // query string
	"www.example.com/data#one=two",                                // query string with #
	"www.example.com/data?one=two&two=http",                       // query string no protocol
	"https://www.example.com/data?one=two&two=http",               // query string two query params, no /
	"https://www.example.com/data/?one=two&two=http",              // query string two query params, with /
	"http://www.example.com/data?one=two&two=https://yup.com?he",  // query string without / and with wierd query param
	"http://www.example.com/data/?one=two&two=https://yup.com?he", // query string with / and with wierd query param
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

func TestTruncateUrlSpecialCasesSuccess(t *testing.T) {
	// Arrange
	t.Parallel()
	assert := require.New(t)
	data := [][]byte{
		GooglePlay,
		ITunesApple,
	}

	// Act
	for _, url := range data {
		result := TruncateUrl(url)
		// Assert
		assert.Equal(url, result)
	}
}

func BenchmarkTrimUrlForScylla(b *testing.B) {
	for _, url := range data2 {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, _, _ = TrimUrlForScylla(url)
		}
	}
}

func TestTrimUrlForScylla(t *testing.T) {
	assert := require.New(t)
	tests := []struct {
		fullUrl  string
		wantUrl  string
		wantHost string
	}{
		{
			fullUrl:  "http://google.com",
			wantUrl:  "https://google.com/",
			wantHost: "google.com",
		},
		{
			fullUrl:  "   http://google.com   ",
			wantUrl:  "https://google.com/",
			wantHost: "google.com",
		},
		{
			fullUrl:  "http://google.com/",
			wantUrl:  "https://google.com/",
			wantHost: "google.com",
		},
		{
			fullUrl:  "https://www.example.com/path?param=value",
			wantUrl:  "https://www.example.com/path/",
			wantHost: "www.example.com",
		},
		{
			fullUrl:  "https://worldtravelling.com/30-stars-we-cant-believe-are-the-same-age/3/?utm_source=Facebook&utm_medium=FB&utm_campaign=DUP GZM_Big4_Vidazoo_CB_Stars The Same Age_P16_RSE - vv6WT WT FB WW An&utm_term=23854019217350509&layout=inf3&vtype=3&fbclid=IwAR3gbeafMqfoDzOPVu2B3P5QEgKtuydi3LmSU4SOft8xa3Akdzo7M0pUtec_aem_th_Aa0DfW8EaIsTtH4kOPKcCwfqRdQUA0TMYlHcRLLLVU1XMA8B43-t-prW7yMcfGw-_MNhLI8vE0TnopF5fjCUJRk4_KDT9WtJ_XXguF0o8qy4Lw",
			wantUrl:  "https://worldtravelling.com/30-stars-we-cant-believe-are-the-same-age/3/",
			wantHost: "worldtravelling.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.fullUrl, func(_ *testing.T) {
			resultUrl, resultHost, err := TrimUrlForScylla(tt.fullUrl)

			assert.NoError(err)
			assert.Equal(tt.wantUrl, resultUrl)
			assert.Equal(tt.wantHost, resultHost)
		})
	}
}

func TestGetDomainFromUrl(t *testing.T) {
	assert := require.New(t)
	testCases := []struct {
		fullUrl        string
		expectedDomain string
		expectedErr    error
	}{
		{"https://www.google.com", "www.google.com", nil},
		{"https://www.example.com/path?param=value", "www.example.com", nil},
	}
	for _, tc := range testCases {
		domain, err := GetDomainFromUrl(tc.fullUrl)
		assert.Equal(tc.expectedDomain, domain)
		assert.Equal(tc.expectedErr, err)
	}
}

func TestUrl(t *testing.T) {
	u := "https://www.the-crossword-solver.com/word/___+major+%28%22the+great+bear%22+constellation%29/"
	parsed, _ := url.Parse(u)

	parsed.Host = "the-crossword-solver.com"
	parsed.RawPath = "/word/___+major+(\"the+great+bear\"+constellation)"
	// + %20

	require.Equal(t, "https://the-crossword-solver.com/word/___+major+%28%22the+great+bear%22+constellation%29/", parsed.String())
}
