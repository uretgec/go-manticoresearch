package manticoresearch

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type UHCOption func(*HttpClient)

func RegisterUHCDefault(name string, timeout int64, reuse bool) UHCOption {
	return func(a *HttpClient) {
		a.name = name
		a.timeout = timeout
		a.reuse = reuse
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

const DefaultUHCName = "MyBakbibuBot"
const DefaultUHCUserAgent = "Mozilla/5.0 (compatible; MyBakbibuBot/1.0; +https://bakbibu.com)"
const DefaultUHCTimeout = 30 // seconds

type HttpClient struct {
	// name // same as User Agent
	name string

	// timeout
	timeout int64 // http client timeout duration

	// client
	reuse bool

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

	return a
}

func (a HttpClient) StatusOK(code int) bool {
	return code == fiber.StatusOK
}

// get request
func (a HttpClient) Get(url string) (code int, body []byte, errs []error) {
	req := fiber.Get(url)
	req.UserAgent(a.name)
	req.Timeout(time.Duration(a.timeout) * time.Second)

	if a.debug {
		req.Debug()
	}

	return req.Bytes()
}

// post request
func (a HttpClient) Post(url string, payload []byte) (code int, body []byte, errs []error) {
	req := fiber.Post(url)
	req.UserAgent(a.name)
	req.ContentType(fiber.MIMETextPlain)

	if a.timeout > 0 {
		req.Timeout(time.Duration(a.timeout) * time.Second)
	}

	req.Body(payload)
	if a.debug {
		req.Debug()
	}

	return req.Bytes()
}

// post json request
func (a HttpClient) PostJSON(url string, payload interface{}) (code int, body []byte, errs []error) {
	req := fiber.Post(url)
	req.UserAgent(a.name)

	if a.timeout > 0 {
		req.Timeout(time.Duration(a.timeout) * time.Second)
	}

	req.JSON(payload)
	if a.debug {
		req.Debug()
	}

	return req.Bytes()
}

// post ndjson request
func (a HttpClient) PostNDJSON(url string, payload []byte) (code int, body []byte, errs []error) {
	req := fiber.Post(url)
	req.UserAgent(a.name)
	req.ContentType("application/x-ndjson")

	if a.timeout > 0 {
		req.Timeout(time.Duration(a.timeout) * time.Second)
	}

	req.Body(payload)
	if a.debug {
		req.Debug()
	}

	return req.Bytes()
}

// put request
func (a HttpClient) Put(url string, payload []byte) (code int, body []byte, errs []error) {
	req := fiber.Put(url)
	req.UserAgent(a.name)

	if a.timeout > 0 {
		req.Timeout(time.Duration(a.timeout) * time.Second)
	}

	req.Body(payload)
	if a.debug {
		req.Debug()
	}

	return req.Bytes()
}

// put json request
func (a HttpClient) PutJSON(url string, payload interface{}) (code int, body []byte, errs []error) {
	req := fiber.Put(url)
	req.UserAgent(a.name)

	if a.timeout > 0 {
		req.Timeout(time.Duration(a.timeout) * time.Second)
	}

	req.JSON(payload)
	if a.debug {
		req.Debug()
	}

	return req.Bytes()
}
