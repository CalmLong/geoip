package main

import (
	"archive/zip"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const srcName = "GeoLite2-Country-CSV"

var (
	pwd, _  = os.Getwd()
	srcV4   = filepath.Join(pwd, srcName, "GeoLite2-Country-Blocks-IPv4.csv")
	srcV6   = filepath.Join(pwd, srcName, "GeoLite2-Country-Blocks-IPv6.csv")
	srcLang = filepath.Join(pwd, srcName, "GeoLite2-Country-Locations-en.csv")
	srcFile = filepath.Join(pwd, srcName+".zip")
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

func getFile() error {
	_ = os.RemoveAll(srcFile)
	fi, err := os.OpenFile(srcFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	val := url.Values{}
	val.Set("license_key", os.Getenv("MAX_MIND_KEY"))
	val.Set("edition_id", "GeoLite2-Country-CSV")
	val.Set("suffix", "zip")
	reqUrl := fmt.Sprintf("https://download.maxmind.com/app/geoip_download?%s", val.Encode())
	request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return err
	}
	resp, err := req.Do(request)
	if err != nil {
		return err
	}
	_, err = io.Copy(fi, resp.Body)
	return resp.Body.Close()
}

func unzip(fileName string) error {
	r, err := zip.OpenReader(fileName)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(srcName, os.ModePerm); err != nil {
		return err
	}
	for _, z := range r.Reader.File {
		r, err := z.Open()
		if err != nil {
			return err
		}
		NewFile, err := os.Create(filepath.Join(pwd, srcName, filepath.Base(z.Name)))
		if err != nil {
			return err
		}
		_, _ = io.Copy(NewFile, r)
		_ = NewFile.Close()
		_ = r.Close()
	}
	return r.Close()
}

func getSrc() {
	if err := getFile(); err != nil {
		log.Fatal(err)
	}
	if err := unzip(srcFile); err != nil {
		log.Fatal(err)
	}
}
