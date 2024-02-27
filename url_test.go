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
	"http://www.example.com/?one=two&two=https://yup.com?he",
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
		{
			fullUrl:  "https://www.zajenata.bg/напуснах-съпруга-си-заради-това-което-той-искаше-да-наÐfbclid=IwAR1z7R9Na9aasrmffWud0nzSCXwsH8TRa1qswcNFA9XtrM3uKvzAGPYkfMU",
			wantUrl:  "https://www.zajenata.bg/напуснах-съпруга-си-заради-това-което-той-искаше-да-наÐfbclid=IwAR1z7R9Na9aasrmffWud0nzSCXwsH8TRa1qswcNFA9XtrM3uKvzAGPYkfMU/",
			wantHost: "www.zajenata.bg",
		},
		{
			fullUrl:  "https://www-lavanguardia-com.cdn.ampproject.org/v/s/www.lavanguardia.com/historiayvida/edad-moderna/20240225/9525175/caravaggio-pintor-homicida-vida-obra-puro-drama.amp.html?amp_gsa=1&amp_js_v=a9&usqp=mq331AQGsAEggAID#amp_tf=De %1$s&aoh=17090467534993&csi=0&referrer=https://www.google.com&ampshare=https://www.lavanguardia.com/historiayvida/edad-moderna/20240225/9525175/caravaggio-pintor-homicida-vida-obra-puro-drama.html",
			wantUrl:  "https://www-lavanguardia-com.cdn.ampproject.org/v/s/www.lavanguardia.com/historiayvida/edad-moderna/20240225/9525175/caravaggio-pintor-homicida-vida-obra-puro-drama.amp.html/",
			wantHost: "www-lavanguardia-com.cdn.ampproject.org",
		},
		{
			fullUrl:  "https://www.football.london/chelsea-fc/news/var-confirm-controversial-moises-caicedo-28701601.amp#amp_tf=From %1$s&aoh=17090476429887&csi=1&referrer=https://www.google.com&ampshare=https://www.football.london/chelsea-fc/news/var-confirm-controversial-moises-caicedo-28701601",
			wantUrl:  "https://www.football.london/chelsea-fc/news/var-confirm-controversial-moises-caicedo-28701601.amp/",
			wantHost: "www.football.london",
		},
		{
			fullUrl:  "https://www.proplanta.de/Agrar-Wetter/Saarbr/ufffdcken-7-Tage-Wettervorhersage.html",
			wantUrl:  "https://www.proplanta.de/Agrar-Wetter/Saarbr/ufffdcken-7-Tage-Wettervorhersage.html/",
			wantHost: "www.proplanta.de",
		},
		{
			fullUrl:  "https://www.oraridiapertura24.it/filiale/muggi/ufffd-comune%20di%20muggi/ufffd%20(servizi%20demografici)-145894j.html",
			wantUrl:  "https://www.oraridiapertura24.it/filiale/muggi/ufffd-comune di muggi/ufffd (servizi demografici)-145894j.html/",
			wantHost: "www.oraridiapertura24.it",
		},
		{
			fullUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล\\ufffd\\ufffd",
			wantUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล\\ufffd\\ufffd/",
			wantHost: "www.goal.com",
		},
		{
			fullUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล\ufffd\ufffd",
			wantUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล\ufffd\ufffd/",
			wantHost: "www.goal.com",
		},
		{
			fullUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล\xf0\x8c\xbc",
			wantUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล/",
			wantHost: "www.goal.com",
		},
		{
			fullUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล\xf8\xa1\xa1\xa1\xa1",
			wantUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล/",
			wantHost: "www.goal.com",
		},
		{
			fullUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล\ufffd\ufffd",
			wantUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล��/",
			wantHost: "www.goal.com",
		},
		{
			fullUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล��",
			wantUrl:  "https://www.goal.com/th/ข่าว/แดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล\ufffd\ufffd/",
			wantHost: "www.goal.com",
		},
		{
			fullUrl:  "https://www.goal.com/th/ข่าว/ʘԊ ԋꙨ ꙩんԽ խแดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล",
			wantUrl:  "https://www.goal.com/th/ข่าว/ʘԊ ԋꙨ ꙩんԽ խแดงเดือดมาตามนัด-แมนฯ-ยูฯ-ดวล/",
			wantHost: "www.goal.com",
		},
		{
			fullUrl:  "https://www.oraridiapertura24.it/filiale/muggi/ufffd-comune%20di%20muggi/ufffd%20(servizi%20demografici)-145894j.html",
			wantUrl:  "https://www.oraridiapertura24.it/filiale/muggi/ufffd-comune di muggi/ufffd (servizi demografici)-145894j.html/",
			wantHost: "www.oraridiapertura24.it",
		},
		{
			fullUrl:  "https://www.deine-tierwelt.de/kleinanzeigen/voegel-c4114/q-rotfl%E3%BCgelsittich/",
			wantUrl:  "https://www.deine-tierwelt.de/kleinanzeigen/voegel-c4114/q-rotflgelsittich/",
			wantHost: "www.deine-tierwelt.de",
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
