// GeoIP generator
//
// Before running this file, the GeoIP database must be downloaded and present.
// To download GeoIP database: https://dev.maxmind.com/geoip/geoip2/geolite2/
// Inside you will find block files for IPv4 and IPv6 and country code mapping.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
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

func getCountryCodeMap() (map[string]string, error) {
	countryCodeReader, err := os.Open(srcLang)
	if err != nil {
		return nil, err
	}
	defer countryCodeReader.Close()
	
	m := make(map[string]string)
	reader := csv.NewReader(countryCodeReader)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, line := range lines[1:] {
		id := line[0]
		countryCode := line[4]
		if len(countryCode) == 0 {
			continue
		}
		m[id] = strings.ToUpper(countryCode)
	}
	return m, nil
}

func getCidrPerCountry(file string, m map[string]string, list map[string][]*router.CIDR) error {
	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileReader.Close()
	
	reader := csv.NewReader(fileReader)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for _, line := range lines[1:] {
		cidrStr := line[0]
		countryID := line[1]
		if countryCode, found := m[countryID]; found {
			if !strings.EqualFold(countryCode, "CN") {
				countryCode = "NCN"
			}
			cidr, err := conf.ParseIP(cidrStr)
			if err != nil {
				return err
			}
			cidrs := append(list[countryCode], cidr)
			list[countryCode] = cidrs
		}
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
	cn := flag.Bool("c", false, "only output cn ips")
	flag.Parse()
	
	cidrList := make(map[string][]*router.CIDR)
	
	if !*cn {
		getSrc()
		ccMap, err := getCountryCodeMap()
		if err != nil {
			fmt.Println("Error reading country code map:", err)
			os.Exit(1)
		}
		
		if err := getCidrPerCountry(srcV4, ccMap, cidrList); err != nil {
			fmt.Println("Error loading IPv4 file:", err)
			os.Exit(1)
		}
		if err := getCidrPerCountry(srcV6, ccMap, cidrList); err != nil {
			fmt.Println("Error loading IPv6 file:", err)
			os.Exit(1)
		}
	} else {
		if err := getCnIPs(cidrList); err != nil {
			fmt.Println("Error getCnIPs: ", err)
			os.Exit(1)
		}
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
	} else {
		fmt.Println("geoip.dat has been generated successfully in the directory.")
	}
	
	_ = os.RemoveAll(filepath.Join(pwd, srcName))
	_ = os.RemoveAll(srcFile)
}