package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	
	"github.com/golang/protobuf/proto"
	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/infra/conf/rule"
)

const countryCN = "CN"
const chinaIP = "https://raw.githubusercontent.com/carrnot/china-ip-list/release/ip.txt"

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

func getPrivateIPs() *routercommon.GeoIP {
	cidr := make([]*routercommon.CIDR, 0, len(privateIPs))
	for _, ip := range privateIPs {
		c, err := rule.ParseIP(ip)
		common.Must(err)
		cidr = append(cidr, c)
	}
	return &routercommon.GeoIP{
		CountryCode: "PRIVATE",
		Cidr:        cidr,
	}
}

func getIPList() []string {
	resp, err := http.Get(chinaIP)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	
	ipList := make([]string, 0)
	
	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		ipList = append(ipList, sc.Text())
	}
	
	return ipList
}

func main() {
	ipList := getIPList()
	cidrList := make(map[string][]*routercommon.CIDR)
	
	for i := 0; i < len(ipList); i++ {
		cidr, err := rule.ParseIP(ipList[i])
		if err != nil {
			continue
		}
		cidrs := append(cidrList[countryCN], cidr)
		cidrList[countryCN] = cidrs
	}
	
	geoIPList := new(routercommon.GeoIPList)
	for cc, cidr := range cidrList {
		geoIPList.Entry = append(geoIPList.Entry, &routercommon.GeoIP{
			CountryCode: cc,
			Cidr:        cidr,
		})
	}
	geoIPList.Entry = append(geoIPList.Entry, getPrivateIPs())
	
	geoIPBytes, err := proto.Marshal(geoIPList)
	if err != nil {
		log.Fatalln("error marshalling geoip list:", err)
	}
	if err = os.WriteFile("geoip.dat", geoIPBytes, 0644); err != nil {
		log.Fatalln("error writing geoip to file:", err)
	}
}
