package ga

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/wabarc/helper"
)

func html() string {
	buf, err := os.ReadFile("testdata/search.html")
	if err != nil {
		return ""
	}
	return string(buf)
}

func TestWayback(t *testing.T) {
	httpClient, mux, server := helper.MockServer()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/search" {
			_, _ = w.Write([]byte(html()))
		}
	})

	uri := "https://example.com"
	u, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}
	wbrc := &Archiver{Client: httpClient}
	got, err := wbrc.Wayback(context.Background(), u)
	if err != nil {
		t.Log(got)
		t.Fatal(err)
	}
}

func TestPlayback(t *testing.T) {
	httpClient, mux, server := helper.MockServer()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(html()))
	})

	uri := "https://example.com"
	u, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}
	wbrc := &Archiver{Client: httpClient}
	got, err := wbrc.Playback(context.Background(), u)
	if err != nil {
		t.Log(got)
		t.Fatal(err)
	}
}
