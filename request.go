package main

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
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

func getFile() error {
	_ = os.RemoveAll(ipSrcFileName)
	fi, err := os.OpenFile(ipSrcFileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	body, err := newHTTPGet(ipSrc)
	if err != nil {
		return err
	}
	_, err = io.Copy(fi, body)
	return err
}
