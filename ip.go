package main

import (
	"bufio"
	"github.com/v2fly/v2ray-core/v4/app/router"
	"github.com/v2fly/v2ray-core/v4/common"
	"github.com/v2fly/v2ray-core/v4/infra/conf/rule"
	"io"
	"log"
)

var ipV4, ipV6 []string

func getPrivateIPs() *router.GeoIP {
	cidr := make([]*router.CIDR, 0, len(privateIPs))
	for _, ip := range privateIPs {
		c, err := rule.ParseIP(ip)
		common.Must(err)
		cidr = append(cidr, c)
	}
	return &router.GeoIP{
		CountryCode: "PRIVATE",
		Cidr:        cidr,
	}
}

func getIPFromUrls(ipUrls []string) []string {
	ipList := make([]string, 0)
	for _, ip := range ipUrls {
		body, err := newHTTPGet(ip)
		if err != nil {
			log.Fatalln("getIPFromUrls err:", err)
		}
		buff := bufio.NewReader(body)
		for {
			s, _, e := buff.ReadLine()
			if e == io.EOF {
				break
			}
			ipList = append(ipList, string(s))
		}
	}
	return ipList
}

func initGeoIP() {
	for _, ip := range ipv4s {
		body, err := newHTTPGet(ip)
		if err != nil {
			log.Fatalln("initGeoIP err:", err)
		}
		buff := bufio.NewReader(body)
		for {
			s, _, e := buff.ReadLine()
			if e == io.EOF {
				break
			}
			ipV4 = append(ipV4, string(s))
		}
	}
	ipV4 = getIPFromUrls(ipv4s)
	ipV6 = getIPFromUrls(ipv6s)
}

func init() {
	initGeoIP()
}
