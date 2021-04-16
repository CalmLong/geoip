package main

const countryCN = "CN"

const (
	ipV41 = "https://raw.githubusercontent.com/gaoyifan/china-operator-ip/ip-lists/china.txt"
	ipV42 = "https://raw.githubusercontent.com/gaoyifan/china-operator-ip/ip-lists/googlecn.txt"
	
	ipV61 = "https://raw.githubusercontent.com/gaoyifan/china-operator-ip/ip-lists/china6.txt"
	ipV62 = "https://raw.githubusercontent.com/gaoyifan/china-operator-ip/ip-lists/googlecn6.txt"
)

var ipv4s = []string{ipV41, ipV42}
var ipv6s = []string{ipV61, ipV62}
