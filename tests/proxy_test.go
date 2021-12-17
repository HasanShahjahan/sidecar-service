package tests

import (
	"github.com/HasanShahjahan/sidecar-service/internal/proxy"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var allowedList = []*regexp.Regexp{
	regexp.MustCompile(`^company$`),
	regexp.MustCompile(strings.Replace(`^company/{id}$`, "{id}", proxy.IdFormat, 1)),
	regexp.MustCompile(`^company/account$`),
	regexp.MustCompile(`^account$`),
	regexp.MustCompile(strings.Replace(`^account/{id}$`, "{id}", proxy.IdFormat, 1)),
	regexp.MustCompile(strings.Replace(`^{id}$`, "{id}", proxy.IdFormat, 1)),
	regexp.MustCompile(strings.Replace(`^account/{id}/user$`, "{id}", proxy.IdFormat, 1)),
	regexp.MustCompile(`^tenant/account/blocked$`),
}

func TestValidator(t *testing.T) {
	var cases = []struct {
		path      string
		expection bool
	}{
		{"company", true},
		{"tenant/sj3co3s4", false},
		{"company/sd45f768", true},
		{"account/acc74850", true},
		{"company/account", true},
		{"acc734340", true},
		{"account/acc234234/user", true},
		{"account/blocked", false},
		{"tenant/account/blocked", true},
		{"tenant/account/acc23849", false},
	}

	proxy := proxy.Proxy{
		Downstream: &proxy.Downstream{
			AllowedList: allowedList,
		},
	}

	for _, tc := range cases {
		require.Equal(t, tc.expection, proxy.ValidatePath(tc.path), "Test is failing! %d equal %s", tc.expection, tc.path)
	}
}

func TestProxyHandler(t *testing.T) {
	var cases = map[string]struct {
		method string
		path   string
		status int
		body   string
		cookie *http.Cookie
	}{
		"/company":          {"GET", "/company", http.StatusOK, "I am Hasan", &http.Cookie{Name: "someName", Value: "SomeValue"}},
		"/company/sd45f768": {"POST", "/company/sd45f768", http.StatusOK, "Done", nil},
		"/company/abc85033": {"PUT", "/company/abc85033", http.StatusBadRequest, "Bad Request!", nil},
		"/account/acc74850": {"GET", "/account/acc74850", http.StatusNotFound, proxy.PathNotFound, nil},
		"/account/acc74850?foo=bar": {"GET", "/account/acc74850?foo=bar&x=y", http.StatusNotFound, proxy.PathNotFound, nil},

	}

	// downstream server
	downstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tCase := cases[r.RequestURI]

		require.Equal(t, tCase.path, r.URL.Path)
		require.Equal(t, tCase.method, r.Method)
		require.NotEmpty(t, r.Header.Get("X-Forwarded-For"))

		if tCase.cookie != nil {
			http.SetCookie(w, tCase.cookie)
		}
		w.WriteHeader(tCase.status)
		w.Write([]byte(tCase.body))
	}))
	defer downstreamServer.Close()

	proxy := proxy.NewProxy(allowedList)

	// sidecar server
	sidecarServer := httptest.NewServer(proxy)
	defer sidecarServer.Close()
	sidecarClient := sidecarServer.Client()

	for _, tCase := range cases {
		getReq, _ := http.NewRequest(tCase.method, sidecarServer.URL + tCase.path, nil)
		getReq.Close = true
		res, err := sidecarClient.Do(getReq)
		if err != nil {
			t.Fatalf("Get: %v", err)
		}

		require.Equal(t, tCase.status, res.StatusCode)
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("reading body: %v", err)
		}

		require.Equal(t, tCase.body, string(body))
		if tCase.cookie != nil {
			require.Equal(t, tCase.cookie.Name, res.Cookies()[0].Name)
			require.Equal(t, tCase.cookie.Value, res.Cookies()[0].Value)
		} else {
			require.Equal(t, 0, len(res.Header["Set-Cookie"]))
		}
	}
}
