package http

import (
	"io"
	"time"
)

type Header map[string]string

type Response struct {
	Status     string
	StatusCode int
	Body       io.Reader
	Header     Header
}

type ClientInterface interface {
	SetHeaders(headers Header)
	RemoveHeader(key string)
	AddHeader(key, value string)
	SetTimeout(d time.Duration)

	Get(url string) (*Response, error)
	Post(url string, body io.Reader) (*Response, error)
	Put(url string, body io.Reader) (*Response, error)
	Patch(url string, body io.Reader) (*Response, error)
	Delete(url string) (*Response, error)
	Download(url string, writer io.Writer) error
}