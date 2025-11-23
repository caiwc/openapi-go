package openapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"log"
	"time"
)

type Signer interface {
	Sign(expires time.Time, request Request) error
}

type ApiKey struct {
	AccessKey, SecretKey string
}

func (a *ApiKey) Sign(expires time.Time, r Request) error {
	log.Printf("first accessKey: %s, secretKey: %s \n", a.AccessKey, a.SecretKey)
	s, err := r.StringToSign()
	if err != nil {
		return err
	}
	log.Printf("second accessKey: %s, secretKey: %s \n", a.AccessKey, a.SecretKey)
	t := fmt.Sprintf("%d", expires.Unix())
	stringToSign := fmt.Sprintf("%s\n%s\n%s", s, t, a.AccessKey)
	mac := hmac.New(sha1.New, []byte(a.SecretKey))
	mac.Write([]byte(stringToSign))
	signature := mac.Sum(nil)
	r.SetSignature(a.AccessKey, fmt.Sprintf("%x", signature), expires)
	log.Printf("third accessKey: %s, secretKey: %s \n", a.AccessKey, a.SecretKey)
	return nil
}
