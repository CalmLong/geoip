package main

import (
	"flag"
	"github.com/golang/protobuf/proto"
	"github.com/v2fly/v2ray-core/v4/infra/conf/rule"
	"io/ioutil"
	"log"
	
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

func ipGeoIP() {
	if err := writer2File(nil, "geoip.txt", ipV4, ipV6); err != nil {
		log.Fatalln("writer2File err: ", err)
	}
}

func clashPGeoIP() {
	v4 := formatIP("  - IP-CIDR,%s", ipV4)
	v6 := formatIP("  - IP-CIDR6,%s", ipV6)
	header := []string{"payload:"}
	if err := writer2File(header, "clashP.yaml", v4, v6); err != nil {
		log.Fatalln("writer2File err: ", err)
	}
}

func v2rayGeoIP() {
	cidrList := make(map[string][]*router.CIDR)
	for _, v4 := range ipV4 {
		cidr, err := rule.ParseIP(v4)
		if err != nil {
			continue
		}
		cidrs := append(cidrList[countryCN], cidr)
		cidrList[countryCN] = cidrs
	}
	for _, v6 := range ipV6 {
		cidr, err := rule.ParseIP(v6)
		if err != nil {
			continue
		}
		cidrs := append(cidrList[countryCN], cidr)
		cidrList[countryCN] = cidrs
	}
	geoIPList := new(router.GeoIPList)
	for cc, cidr := range cidrList {
		geoIPList.Entry = append(geoIPList.Entry, &router.GeoIP{
			CountryCode: cc,
			Cidr:        cidr,
		})
	}
	geoIPList.Entry = append(geoIPList.Entry, getPrivateIPs())
	
	geoIPBytes, err := proto.Marshal(geoIPList)
	if err != nil {
		log.Fatalln("error marshalling geoip list:", err)
	}
	if err := ioutil.WriteFile("geoip.dat", geoIPBytes, 0644); err != nil {
		log.Fatalln("error writing geoip to file:", err)
	}
}

func command() {
	fF := flag.String("F", "v2ray", "")
	flag.Parse()
	switch *fF {
	case "clashP":
		clashPGeoIP()
	case "ip":
		ipGeoIP()
	case "v2ray":
		fallthrough
	default:
		v2rayGeoIP()
	}
}

func main() {
	command()
	log.Println("geoip has been generated successfully in the directory.")
}
