package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var req = http.Client{Transport: Transport()}

func Transport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		DialContext: (&net.Dialer{
			Timeout: 180 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 60 * time.Second,
		ForceAttemptHTTP2:   true,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 50,
	}
}

func newHTTPGet(uri string) (io.Reader, error) {
	log.Println(uri)
	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := req.Do(request)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
