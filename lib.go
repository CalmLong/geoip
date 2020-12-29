package main

const countryCN = "CN"

const (
	ipSrcFileName = "iplist-master.zip"
	ipSrcName     = "iplist-master"
	ipSrc         = "https://github.com/metowolf/iplist/archive/master.zip"
	ipData        = "iplist-master/data/country"
	ipDataCN      = "iplist-master/data/special/china.txt"
)

const (
	ipV61 = "https://raw.githubusercontent.com/gaoyifan/china-operator-ip/ip-lists/china6.txt"
)

var ipv6s = []string{ipV61}

var pwd string
