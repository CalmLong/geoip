package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"v2ray.com/core/app/router"
	"v2ray.com/core/infra/conf"
)

const cnIPUrl = "https://raw.githubusercontent.com/metowolf/iplist/master/data/special/china.txt"

func getCnIPs(list map[string][]*router.CIDR) error {
	request, err := http.NewRequest(http.MethodGet, cnIPUrl, nil)
	if err != nil {
		return err
	}
	resp, err := req.Do(request)
	if err != nil {
		return err
	}
	ips := bufio.NewReader(resp.Body)
	for {
		ip, _, e := ips.ReadLine()
		if e == io.EOF {
			break
		}
		cidr, err := conf.ParseIP(fmt.Sprintf("%s", ip))
		if err != nil {
			return err
		}
		cidrs := append(list["CN"], cidr)
		list["CN"] = cidrs
	}
	return nil
}