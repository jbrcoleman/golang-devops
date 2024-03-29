package api

import (
	"io"
	"net/http"
)

type ClientIface interface {
	Get(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Options struct {
	Password string
	LoginURL string
}

type APIIface interface {
	DoGetRequest(requestURL string) (Response, error)
}

type api struct {
	Options Options
	Client  ClientIface
}

func New(options Options) APIIface {
	return api{
		Options: options,
		Client: &http.Client{
			Transport: &MYJWTTransport{
				transport:  http.DefaultTransport,
				password:   options.Password,
				loginURL:   options.LoginURL,
				HTTPClient: &http.Client{},
			},
		},
	}
}
