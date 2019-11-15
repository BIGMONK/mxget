package sreq

import (
	"errors"
	"net"
	"net/http"
	"net/http/cookiejar"
	stdurl "net/url"
	"sync"
	"time"

	"golang.org/x/net/publicsuffix"
)

var std = New(nil)

type (
	// Client defines a sreq client and will be reused for per request.
	Client struct {
		// RawClient specifies an HTTP client for sending HTTP requests.
		RawClient *http.Client

		// RequestOptions specifies request options that sreq uses for per HTTP request by default.
		RequestOptions []RequestOption

		mux sync.RWMutex
	}
)

// DefaultHTTPClient returns an HTTP client that sreq uses by default.
func DefaultHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	timeout := 120 * time.Second
	return &http.Client{
		Transport: transport,
		Jar:       jar,
		Timeout:   timeout,
	}
}

// New allows you to customize a sreq client with an HTTP client.
// If the transport or timeout of the HTTP client not specified, sreq would use defaults.
func New(httpClient *http.Client) *Client {
	rawClient := DefaultHTTPClient()
	if httpClient != nil {
		if httpClient.Transport != nil {
			rawClient.Transport = httpClient.Transport
		}
		if httpClient.Timeout > 0 {
			rawClient.Timeout = httpClient.Timeout
		}
		rawClient.CheckRedirect = httpClient.CheckRedirect
		rawClient.Jar = httpClient.Jar
	}

	return &Client{
		RawClient: rawClient,
	}
}

// SetDefaultRequestOpts sets default request options for per HTTP request.
func SetDefaultRequestOpts(opts ...RequestOption) {
	std.SetDefaultRequestOpts(opts...)
}

// SetDefaultRequestOpts sets default request options for per HTTP request.
func (c *Client) SetDefaultRequestOpts(opts ...RequestOption) {
	c.mux.Lock()
	c.RequestOptions = opts
	c.mux.Unlock()
}

// AddDefaultRequestOpts appends default request options for per HTTP request.
func AddDefaultRequestOpts(opts ...RequestOption) {
	std.AddDefaultRequestOpts(opts...)
}

// AddDefaultRequestOpts appends default request options for per HTTP request.
func (c *Client) AddDefaultRequestOpts(opts ...RequestOption) {
	c.mux.Lock()
	c.RequestOptions = append(c.RequestOptions, opts...)
	c.mux.Unlock()
}

// ClearDefaultRequestOpts clears default request options for per HTTP request.
func ClearDefaultRequestOpts() {
	std.ClearDefaultRequestOpts()
}

// ClearDefaultRequestOpts clears default request options for per HTTP request.
func (c *Client) ClearDefaultRequestOpts() {
	c.mux.Lock()
	c.RequestOptions = nil
	c.mux.Unlock()
}

// FilterCookies returns the cookies to send in a request for the given URL.
func FilterCookies(url string) ([]*http.Cookie, error) {
	return std.FilterCookies(url)
}

// FilterCookies returns the cookies to send in a request for the given URL.
func (c *Client) FilterCookies(url string) ([]*http.Cookie, error) {
	if c.RawClient.Jar == nil {
		return nil, errors.New("sreq: nil cookie jar")
	}

	u, err := stdurl.Parse(url)
	if err != nil {
		return nil, err
	}
	cookies := c.RawClient.Jar.Cookies(u)
	if len(cookies) == 0 {
		return nil, errors.New("sreq: cookies for the given URL not present")
	}

	return cookies, nil
}

// FilterCookie returns the named cookie to send in a request for the given URL.
func FilterCookie(url string, name string) (*http.Cookie, error) {
	return std.FilterCookie(url, name)
}

// FilterCookie returns the named cookie to send in a request for the given URL.
func (c *Client) FilterCookie(url string, name string) (*http.Cookie, error) {
	cookies, err := c.FilterCookies(url)
	if err != nil {
		return nil, err
	}

	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie, nil
		}
	}

	return nil, errors.New("sreq: named cookie for the given URL not present")
}

// Do sends a raw HTTP request and returns its response.
func Do(rawRequest *http.Request) *Response {
	return std.Do(rawRequest)
}

// Do sends a raw HTTP request and returns its response.
func (c *Client) Do(rawRequest *http.Request) *Response {
	rawResponse, err := c.RawClient.Do(rawRequest)
	return &Response{
		RawResponse: rawResponse,
		Err:         err,
	}
}
