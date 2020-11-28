package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/golang/protobuf/proto"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/infra/conf"
)

var privateIPs = []string{
	"0.0.0.0/8",
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

func getCidrPerCountry(country []string, list map[string][]*router.CIDR) error {
	for _, code := range country {
		countryStr := filepath.Join(pwd, ipData, code+".txt")
		if strings.EqualFold(code, "CN") {
			countryStr = filepath.Join(pwd, ipDataCN)
		}
		countryIps, err := os.Open(countryStr)
		if err != nil {
			return err
		}
		cidr := bufio.NewReader(countryIps)
		for {
			cidrStr, _, err := cidr.ReadLine()
			if err == io.EOF {
				break
			}
			cidr, err := conf.ParseIP(string(cidrStr))
			if err != nil {
				return err
			}
			cidrs := append(list[code], cidr)
			list[code] = cidrs
		}
		countryIps.Close()
	}
	return nil
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
	var err error
	pwd, err = os.Getwd()
	if err != nil {
		log.Fatal("get pwd: ", err)
	}
	
	cidrList := make(map[string][]*router.CIDR)
	
	if err := getFile(); err != nil {
		log.Fatal("getFile: ", err)
	}
	if err := unzip(); err != nil {
		log.Fatal("unzip: ", err)
	}
	
	ccMap, err := getCountryList()
	if err != nil {
		log.Fatal("getCountryList: ", err)
	}
	
	if err := getCidrPerCountry(ccMap, cidrList); err != nil {
		log.Fatal("getCidrPerCountry: ", err)
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
		fmt.Println("Error marshalling geoip list:", err)
		os.Exit(1)
	}
	
	if err := ioutil.WriteFile("geoip.dat", geoIPBytes, 0644); err != nil {
		fmt.Println("Error writing geoip to file:", err)
		os.Exit(1)
	}
	
	fmt.Println("geoip.dat has been generated successfully in the directory.")
	_ = os.RemoveAll(filepath.Join(pwd, ipSrcName))
	_ = os.RemoveAll(filepath.Join(pwd, ipSrcFileName))
}