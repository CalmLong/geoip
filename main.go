package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/v2fly/v2ray-core/v4/common"
	"github.com/v2fly/v2ray-core/v4/common/net"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
	
	"github.com/v2fly/v2ray-core/v4/app/router"
)

var privateIPs = []string{
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.0.0.0/24",
	"192.0.2.0/24",
	"192.88.99.0/24",
	"192.168.0.0/16",
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"224.0.0.0/4",
	"240.0.0.0/4",
	"255.255.255.255/32",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
}

func getChinaIPsFromRawUrl(ipUrls []string, list map[string][]*router.CIDR) {
	for _, ip := range ipUrls {
		body, err := newHTTPGet(ip)
		if err != nil {
			log.Fatalln("getChinaIPsFromRawUrl err:", err)
		}
		buff := bufio.NewReader(body)
		for {
			s, _, e := buff.ReadLine()
			if e == io.EOF {
				break
			}
			cidr, err := conf.ParseIP(string(s))
			if err != nil {
				continue
			}
			cidrs := append(list[countryCN], cidr)
			list[countryCN] = cidrs
		}
	}
}

func pickWriter(header []string, name string, items []string) error {
	buff := bytes.NewBuffer([]byte{})
	for _, h := range header {
		buff.WriteString(h)
	}
	for _, item := range items {
		buff.WriteString(item)
	}
	return ioutil.WriteFile(name, buff.Bytes(), os.ModePerm)
}

func getChinaCidr(list map[string][]*router.CIDR) {
	getChinaIPsFromRawUrl(ipv4s, list)
	getChinaIPsFromRawUrl(ipv6s, list)
}

func getPrivateIPs() *router.GeoIP {
	cidr := make([]*router.CIDR, 0, len(privateIPs))
	for _, ip := range privateIPs {
		c, err := conf.ParseIP(ip)
		common.Must(err)
		cidr = append(cidr, c)
	}
	return &router.GeoIP{
		CountryCode: "PRIVATE",
		Cidr:        cidr,
	}
}

func main() {
	fF := flag.String("F", "", "")
	flag.Parse()
	
	cidrList := make(map[string][]*router.CIDR)
	getChinaCidr(cidrList)
	
	geoIPList := new(router.GeoIPList)
	for cc, cidr := range cidrList {
		geoIPList.Entry = append(geoIPList.Entry, &router.GeoIP{
			CountryCode: cc,
			Cidr:        cidr,
		})
	}
	
	t := time.Now().Format("2006-01-02 15:04:05")
	
	switch *fF {
	case "clash":
		ips := make([]string, 0)
		const rule = "  - IP-CIDR,%s/%d\n"
		
		for _, ip := range getPrivateIPs().GetCidr() {
			ipStr := fmt.Sprintf(rule, net.IPAddress(ip.GetIp()).String(), ip.Prefix)
			ips = append(ips, ipStr)
		}
		
		for _, geo := range geoIPList.GetEntry() {
			if geo.GetCountryCode() == countryCN {
				for _, ip := range geo.GetCidr() {
					ipStr := fmt.Sprintf(rule, net.IPAddress(ip.GetIp()).String(), ip.Prefix)
					ips = append(ips, ipStr)
				}
			}
		}
		
		header := []string{"# TIME: ", t, "\n", "payload:", "\n"}
		if err := pickWriter(header, "geoip.yaml", ips); err != nil {
			log.Fatalln("pickWriter err: ", err)
		}
	default:
		geoIPList.Entry = append(geoIPList.Entry, getPrivateIPs())
		
		geoIPBytes, err := proto.Marshal(geoIPList)
		if err != nil {
			log.Fatalln("error marshalling geoip list:", err)
		}
		if err := ioutil.WriteFile("geoip.dat", geoIPBytes, 0644); err != nil {
			log.Fatalln("error writing geoip to file:", err)
		}
	}
	
	log.Println("geoip has been generated successfully in the directory.")
}
