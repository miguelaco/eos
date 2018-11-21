package common

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const AuthCookieName = "dcos-acs-auth-cookie"

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
	return &HttpClient{client: client}
}

type HttpClient struct {
	client *http.Client
	Token  string
}

func (c *HttpClient) ResetReferer(resetReferer bool) {
	c.client.CheckRedirect = nil
	if resetReferer {
		c.client.CheckRedirect = resetRefererFunc
	}
}

func (c *HttpClient) Verbose(verbose bool) {
	c.client.Transport = http.DefaultTransport
	if verbose {
		c.client.Transport = logRedirects{}
	}
}

func (c *HttpClient) GetCookie(addr string, name string) (result string, err error) {
	rootUrl, _ := url.Parse(addr)

	for _, cookie := range c.client.Jar.Cookies(rootUrl) {
		if cookie.Name == name {
			result = cookie.Value
			return
		}
	}

	err = errors.New("Cookie " + name + " not found")
	return
}

func (c *HttpClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	if c.Token != "" {
		req.AddCookie(&http.Cookie{Name: AuthCookieName, Value: c.Token})
	}

	return c.client.Do(req)
}

func (c *HttpClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if c.Token != "" {
		req.AddCookie(&http.Cookie{Name: AuthCookieName, Value: c.Token})
	}

	return c.client.Do(req)
}
