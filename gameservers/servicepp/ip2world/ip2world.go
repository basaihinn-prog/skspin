package main

import (
	"fmt"
	"net/http"
	"serve/comm/ut"
	"strconv"
	"strings"

	ip2worldpub "serve/servicepp/ip2world/pub"
)

var ip2worldurl = "https://gw.dataimpulse.com:777/api/list?format=hostname:port&quantity=100&type=sticky"

//var ip2worldurl = "https://485189d8e612a951392d:5853122d271e64a8gw.dataimpulse.com:777/api/list?format=hostname:port&quantity=100&type=sticky"

type Endpoint = ip2worldpub.Endpoint

func getEndpoints(regions string) (endpoints []Endpoint, err error) {
	//ul := fmt.Sprintf(ip2worldurl, regions)

	req, _ := http.NewRequest("GET", ip2worldurl, nil)
	fmt.Println("req url", ip2worldurl)
	body, _, err := ut.DoHttpReq(http.DefaultClient, req)
	if err != nil {
		return
	}

	// 使用 strings.Split 按换行符分割字符串
	parts := strings.Split(string(body), "\n")
	//parts := strings.Split("gw.dataimpulse.com:10000\ngw.dataimpulse.com:10001\ngw.dataimpulse.com:10002\ngw.dataimpulse.com:10003\ngw.dataimpulse.com:10004\ngw.dataimpulse.com:10005\ngw.dataimpulse.com:10006\ngw.dataimpulse.com:10007\ngw.dataimpulse.com:10008\ngw.dataimpulse.com:10009\ngw.dataimpulse.com:10010\ngw.dataimpulse.com:10011\ngw.dataimpulse.com:10012\ngw.dataimpulse.com:10013\ngw.dataimpulse.com:10014\ngw.dataimpulse.com:10015\ngw.dataimpulse.com:10016\ngw.dataimpulse.com:10017\ngw.dataimpulse.com:10018\ngw.dataimpulse.com:10019\ngw.dataimpulse.com:10020\ngw.dataimpulse.com:10021\ngw.dataimpulse.com:10022\ngw.dataimpulse.com:10023\ngw.dataimpulse.com:10024\ngw.dataimpulse.com:10025\ngw.dataimpulse.com:10026\ngw.dataimpulse.com:10027\ngw.dataimpulse.com:10028\ngw.dataimpulse.com:10029\ngw.dataimpulse.com:10030\ngw.dataimpulse.com:10031\ngw.dataimpulse.com:10032\ngw.dataimpulse.com:10033\ngw.dataimpulse.com:10034\ngw.dataimpulse.com:10035\ngw.dataimpulse.com:10036\ngw.dataimpulse.com:10037\ngw.dataimpulse.com:10038\ngw.dataimpulse.com:10039\ngw.dataimpulse.com:10040\ngw.dataimpulse.com:10041\ngw.dataimpulse.com:10042\ngw.dataimpulse.com:10043\ngw.dataimpulse.com:10044\ngw.dataimpulse.com:10045\ngw.dataimpulse.com:10046\ngw.dataimpulse.com:10047\ngw.dataimpulse.com:10048\ngw.dataimpulse.com:10049\ngw.dataimpulse.com:10050\ngw.dataimpulse.com:10051\ngw.dataimpulse.com:10052\ngw.dataimpulse.com:10053\ngw.dataimpulse.com:10054\ngw.dataimpulse.com:10055\ngw.dataimpulse.com:10056\ngw.dataimpulse.com:10057\ngw.dataimpulse.com:10058\ngw.dataimpulse.com:10059\ngw.dataimpulse.com:10060\ngw.dataimpulse.com:10061\ngw.dataimpulse.com:10062\ngw.dataimpulse.com:10063\ngw.dataimpulse.com:10064\ngw.dataimpulse.com:10065\ngw.dataimpulse.com:10066\ngw.dataimpulse.com:10067\ngw.dataimpulse.com:10068\ngw.dataimpulse.com:10069\ngw.dataimpulse.com:10070\ngw.dataimpulse.com:10071\ngw.dataimpulse.com:10072\ngw.dataimpulse.com:10073\ngw.dataimpulse.com:10074\ngw.dataimpulse.com:10075\ngw.dataimpulse.com:10076\ngw.dataimpulse.com:10077\ngw.dataimpulse.com:10078\ngw.dataimpulse.com:10079\ngw.dataimpulse.com:10080\ngw.dataimpulse.com:10081\ngw.dataimpulse.com:10082\ngw.dataimpulse.com:10083\ngw.dataimpulse.com:10084\ngw.dataimpulse.com:10085\ngw.dataimpulse.com:10086\ngw.dataimpulse.com:10087\ngw.dataimpulse.com:10088\ngw.dataimpulse.com:10089\ngw.dataimpulse.com:10090\ngw.dataimpulse.com:10091\ngw.dataimpulse.com:10092\ngw.dataimpulse.com:10093\ngw.dataimpulse.com:10094\ngw.dataimpulse.com:10095\ngw.dataimpulse.com:10096\ngw.dataimpulse.com:10097\ngw.dataimpulse.com:10098\ngw.dataimpulse.com:10099", "\n")

	// 输出分割后的结果
	for _, part := range parts {
		fmt.Println(part)
		url := strings.Split(part, ":")
		if len(url) < 2 {
			fmt.Println(url)
			continue
		}
		fmt.Println(url, len(url)-1, url[len(url)-1])
		port, err := strconv.ParseInt(url[len(url)-1], 10, 64)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, Endpoint{IP: url[0], Port: port})
	}
	return
}
