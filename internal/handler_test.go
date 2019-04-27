package internal_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davyj0nes/go-proxy/internal"
	"github.com/sirupsen/logrus"
)

func TestProxyHandler(t *testing.T) {
	want := "this call was relayed by the reverse proxy"
	backendSrv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, want)
		}),
	)
	defer backendSrv.Close()

	stubLogger := &logrus.Logger{}
	handler := internal.NewHandler(stubLogger, backendSrv.URL)

	proxy := httptest.NewServer(handler)
	resp, err := http.Get(proxy.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if want != string(body) {
		t.Fatalf("want: (%s), got: (%s)", want, string(body))
	}
}

func BenchmarkProxyHandler(b *testing.B) {
	backendSrv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "this call was relayed by the reverse proxy")
		}))
	defer backendSrv.Close()

	stubLogger := &logrus.Logger{}
	handler := internal.NewHandler(stubLogger, backendSrv.URL)
	proxy := httptest.NewServer(handler)

	for n := 0; n < b.N; n++ {
		resp, err := http.Get(proxy.URL)
		if err != nil {
			b.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("expected: 200, got: %d", resp.StatusCode)
		}
	}
}
