package manticoresearch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

type UHCOption func(*HttpClient)

func RegisterUHCDefault(name string, timeout int64) UHCOption {
	return func(a *HttpClient) {
		a.name = name
		a.timeout = timeout
	}
}

func RegisterUHCDebugMode(enabled bool) UHCOption {
	return func(a *HttpClient) {
		a.debug = enabled
	}
}

func RegisterUHCDebugEnable() UHCOption {
	return func(a *HttpClient) {
		a.debug = true
	}
}

const DefaultUHCName = "GoManticoreSearchBot"
const DefaultUHCUserAgent = "Mozilla/5.0 (compatible; GoManticoreSearchBot/1.0; +https://manticoresearch.com)"
const DefaultUHCTimeout = 10 // seconds

type HttpClient struct {
	// name // same as User Agent
	name string

	// timeout
	timeout int64 // http client timeout duration

	// client
	client *http.Client

	// debug mode
	debug bool
}

func New(options ...UHCOption) *HttpClient {
	a := &HttpClient{}

	for _, opt := range options {
		opt(a)
	}

	if a.name == "" {
		a.name = DefaultUHCUserAgent
	}

	if a.timeout == 0 {
		a.timeout = DefaultUHCTimeout
	}

	// update connection pool -> from 2 to 100
	// Reuse http client more requests/responses
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	// Client
	a.client = &http.Client{
		Timeout:   time.Duration(a.timeout) * time.Second,
		Transport: t,
	}

	return a
}

func (a HttpClient) StatusOK(code int) bool {
	return code == http.StatusOK
}

// get request
func (a HttpClient) Get(url string) (code int, body []byte, err error) {
	return a._request(http.MethodGet, url, nil, nil, false)
}

// post request
func (a HttpClient) Post(url string, payload []byte) (code int, body []byte, err error) {
	headers := map[string]string{
		"Content-Type": "text/plain",
	}
	return a._request(http.MethodPost, url, headers, bytes.NewBuffer(payload), false)
}

// post json request
func (a *HttpClient) PostJSON(url string, payload []byte) (code int, body []byte, err error) {
	// headers
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	return a._request(http.MethodPost, url, headers, bytes.NewBuffer(payload), false)
}

// post ndjson request
func (a HttpClient) PostNDJSON(url string, payload []byte) (code int, body []byte, err error) {
	headers := map[string]string{
		"Content-Type": "application/x-ndjson",
	}
	return a._request(http.MethodPost, url, headers, bytes.NewBuffer(payload), false)
}

// put request
func (a HttpClient) Put(url string, payload []byte) (code int, body []byte, err error) {
	return a._request(http.MethodPut, url, nil, bytes.NewBuffer(payload), false)
}

// put json request
func (a HttpClient) PutJSON(url string, payload []byte) (code int, body []byte, err error) {
	return a._request(http.MethodPut, url, nil, bytes.NewBuffer(payload), false)
}

// Private
func (a *HttpClient) _request(method, url string, headers map[string]string, payload io.Reader, statusOnly bool) (code int, body []byte, err error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return 0, nil, err
	}

	// Add Header
	req.Header.Add("User-Agent", a.name)
	if len(headers) > 0 {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	// debug
	if a.debug {
		reqDump, _ := httputil.DumpRequestOut(req, true)
		fmt.Printf("----------------------------\nREQUEST:\n%s\n", string(reqDump))
		reqDump = nil
	}

	// make request
	resp, err := a.client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	// Close the connection to reuse it
	defer resp.Body.Close()

	if !statusOnly {
		body, err = io.ReadAll(resp.Body)
	}

	if a.debug {
		respDump, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("RESPONSE:\n%s\n", string(respDump))
		respDump = nil
	}

	return resp.StatusCode, body, err
}
