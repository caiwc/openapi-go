package openapi

import (
	"log"
	"net/http"
	"time"
)

type TransportOption func(*transport)

type transport struct {
	rt  http.RoundTripper
	sig *ApiKey // 直接存储结构体值，避免接口值拷贝问题
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("first sig info: %v, accessKey: %s, secretKey: %s\n", t.sig, t.sig.AccessKey, t.sig.SecretKey)
	if err := t.sig.Sign(time.Now(), request{req}); err != nil {
		return nil, err
	}
	log.Printf("full url: %s, accessKey: %s, secretKey: %s\n", req.URL.String(), t.sig.AccessKey, t.sig.SecretKey)
	return t.rt.RoundTrip(req)
}

func NewTransport(accessKey, secretKey string, options ...TransportOption) (http.RoundTripper, error) {
	t := &transport{
		rt: http.DefaultTransport,
		sig: &ApiKey{ // 直接使用值类型，避免接口值拷贝问题
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
