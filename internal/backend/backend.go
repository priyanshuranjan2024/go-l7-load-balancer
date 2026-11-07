package backend

import (
	"fmt"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	URL   *url.URL
	Proxy *httputil.ReverseProxy
}

func New(rawURL string) (*Backend, error) {
	target, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if target.Scheme == "" || target.Host == "" {
		return nil, fmt.Errorf("invalid backend URL: %s", rawURL)
	}

	return &Backend{
		URL:   target,
		Proxy: httputil.NewSingleHostReverseProxy(target),
	}, nil
}
