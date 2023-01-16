package dns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

type DNSResponse struct {
	Status   int  `json:"Status"`
	TC       bool `json:"TC"`
	RD       bool `json:"RD"`
	RA       bool `json:"RA"`
	AD       bool `json:"AD"`
	CD       bool `json:"CD"`
	Question []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	} `json:"Question"`
	Answer []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		TTL  int    `json:"TTL"`
		Data string `json:"data"`
	} `json:"Answer"`
}

func Doh() {
	// 指定上游 DoH 服务器的 URL
	url := "https://cloudflare-dns.com/dns-query?name=example.com&type=A"

	// 发送 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 解析 JSON 数据
	var dnsResp DNSResponse
	json.Unmarshal(data, &dnsResp)

	// 遍历答案中的数据，获取 IP 信息
	for _, answer := range dnsResp.Answer {
		if answer.Type == 1 {
			ip := net.ParseIP(answer.Data)
			if ip != nil {
				fmt.Println(ip)
			}
		}
	}
}
