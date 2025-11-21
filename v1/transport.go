package openapi

import (
	"log"
	"net/http"
	"time"
)

type TransportOption func(*transport)

type transport struct {
	rt  http.RoundTripper
	sig ApiKey
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.sig.Sign(time.Now(), request{req}); err != nil {
		return nil, err
	}
	log.Printf("full url: %s , accessKey: %s, secretKey: %s \n", req.URL.String(), t.sig.AccessKey, t.sig.SecretKey)
	return t.rt.RoundTrip(req)
}

func NewTransport(accessKey, secretKey string, options ...TransportOption) (http.RoundTripper, error) {
	t := &transport{
		rt: http.DefaultTransport,
		sig: ApiKey{
			AccessKey: accessKey,
			SecretKey: secretKey,
		},
	}
	for _, option := range options {
		option(t)
	}
	return t, nil
}

func RoundTripper(rt http.RoundTripper) TransportOption {
	return func(t *transport) {
		if rt != nil {
			t.rt = rt
		}
	}
}
