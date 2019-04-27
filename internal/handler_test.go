package internal_test

import (
	"fmt"
	"github.com/davyj0nes/go-proxy/internal"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyHandler(t *testing.T) {
	want := "this call was relayed by the reverse proxy"
	backendSrv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))
	defer backendSrv.Close()

	handler := internal.NewHandler(backendSrv.URL)

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
