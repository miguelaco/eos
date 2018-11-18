package common

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type logRedirects struct {
	Transport http.RoundTripper
}

func (l logRedirects) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t := l.Transport
	if t == nil {
		t = http.DefaultTransport
	}
	resp, err = t.RoundTrip(req)
	if err != nil {
		return
	}

	fmt.Println(req.Method, "for", req.URL, "status", resp.StatusCode)

	return
}

func resetRefererFunc(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errors.New("stopped after 10 redirects")
	}

	req.Header.Set("Referer", "")

	return nil
}

func NewHttpClient() *HttpClient {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}

	return &HttpClient{Client: client}
}

type HttpClient struct {
	*http.Client
}

func (c *HttpClient) ResetReferer(resetReferer bool) {
	c.CheckRedirect = nil
	if resetReferer {
		c.CheckRedirect = resetRefererFunc
	}
}

func (c *HttpClient) Verbose(verbose bool) {
	c.Transport = http.DefaultTransport
	if verbose {
		c.Transport = logRedirects{}
	}
}

func (c *HttpClient) GetCookie(addr string, name string) (string, error) {
	rootUrl, _ := url.Parse(addr)

	for _, cookie := range c.Jar.Cookies(rootUrl) {
		if cookie.Name == name {
			return cookie.Value, nil
		}
	}

	return "", errors.New("Cookie " + name + " not found")
}
