package ga

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wabarc/logger"
	"github.com/wabarc/proxier"
)

type Archiver struct {
	Client *http.Client
}

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36"

var (
	base = "https://ghostarchive.org"

	endpoint = struct {
		wayback  string
		playback string
	}{
		wayback:  base + "/archive2",
		playback: base + "/search",
	}
)

func init() {
	debug := os.Getenv("DEBUG")
	if debug == "true" || debug == "1" || debug == "on" {
		logger.EnableDebug()
	}
}

// Wayback is the handle of saving webpages to archive.org
func (wbrc *Archiver) Wayback(ctx context.Context, u *url.URL) (result string, err error) {
	wbrc.Client = proxier.NewClient(wbrc.Client).Client

	result, err = wbrc.archive(ctx, u)
	if err != nil {
		return
	}
	return
}

// Playback handle searching archived webpages from archive.is
func (wbrc *Archiver) Playback(ctx context.Context, u *url.URL) (result string, err error) {
	wbrc.Client = proxier.NewClient(wbrc.Client).Client

	result, err = wbrc.search(ctx, u)
	if err != nil {
		return
	}
	return
}

func (wbrc *Archiver) archive(ctx context.Context, u *url.URL) (string, error) {
	uri := u.String()
	data := url.Values{
		"archive": {uri},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.wayback, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := wbrc.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var loc string
	loc = resp.Header.Get("Content-Location")
	if len(loc) > 0 {
		return loc, nil
	}

	loc = resp.Header.Get("Location")
	if len(loc) > 0 {
		return loc, nil
	}

	loc = resp.Request.URL.String()
	if len(loc) > 0 {
		return loc, nil
	}

	// loc, err = wbrc.latest(ctx, u)
	// if err != nil {
	// 	loc = base + uri
	// }

	return loc, nil
}

func (wbrc *Archiver) search(ctx context.Context, u *url.URL) (string, error) {
	return wbrc.latest(ctx, u)
}

func (wbrc *Archiver) latest(ctx context.Context, u *url.URL) (string, error) {
	// https://ghostarchive/search?term=https://example.org

	uri := u.String()
	res := fmt.Sprintf("%s/search?term=%s", base, uri)
	data := url.Values{
		"term": {uri},
	}
	url := endpoint.playback + "?" + data.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err := wbrc.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	slug, ok := doc.Find(`#bodyContent > table > tbody > tr:nth-child(2) > td:nth-child(2) > a`).Attr("href")
	if ok {
		return fmt.Sprintf("%s%s", base, slug), nil
	}

	return res, fmt.Errorf("Not found")
}
