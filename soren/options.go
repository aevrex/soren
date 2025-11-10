package soren

import (
	"crypto/tls"
	"net/http"
)

type Option func(*Client)

func WithInsecureTLS() Option {
	return func(c *Client) {
		tr, ok := c.httpClient.Transport.(*http.Transport)
		if !ok || tr == nil {
			tr = http.DefaultTransport.(*http.Transport).Clone()
		}

		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		}
		tr.TLSClientConfig.InsecureSkipVerify = true

		c.httpClient.Transport = tr
	}
}
