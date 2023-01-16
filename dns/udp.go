package dns

import (
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/MoeQing-Network/MoeDNS/utils"
	"github.com/miekg/dns"
)

func Start() {
	// 监听本地端口 53
	serverAddr, err := net.ResolveUDPAddr("udp", ":53")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer serverConn.Close()

	// 上游 DNS 服务器
	upstreamAddr, err := net.ResolveUDPAddr("udp", "1.0.0.1:53")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 处理请求
	sem := make(chan struct{}, 100)

	for {
		sem <- struct{}{}
		go handleRequest(serverConn, upstreamAddr, sem)
	}
}

func handleRequest(serverConn *net.UDPConn, upstreamAddr *net.UDPAddr, sem chan struct{}) {
	defer func() {
		<-sem
	}()
	// 读取请求数据
	request := make([]byte, 1024)
	n, addr, err := serverConn.ReadFromUDP(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 打印请求的域名
	fmt.Println("Request for", string(request[:n]))

	// 转发请求到上游 DNS 服务器
	upstreamConn, err := net.DialUDP("udp", nil, upstreamAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer upstreamConn.Close()

	// 将请求数据发送到上游 DNS 服务器
	_, err = upstreamConn.Write(request[:n])
	if err != nil {
		fmt.Println(err)
		return
	}

	// 读取上游 DNS 服务器的响应
	response := make([]byte, 1024)
	n, err = upstreamConn.Read(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := new(dns.Msg)
	err = msg.Unpack(response[:n])
	if err != nil {
		fmt.Println(err)
		return
	}

	// 确定响应中包含的 IP 数量
	var ipCount int
	var matchPrefix bool
	for _, answer := range msg.Answer {
		if a, ok := answer.(*dns.A); ok {
			ip := a.A
			fmt.Println(ip.String())
			ipCount++
		}
		if a, ok := answer.(*dns.AAAA); ok {
			ip := a.AAAA
			ip_net, _ := netip.ParseAddr(ip.String())
			if utils.FindPrefix(ip_net) {
				matchPrefix = true
			}
			ipCount++
		}

	}
	fmt.Println("Total IPs:", ipCount)

	if matchPrefix {
		for i := 0; i < len(msg.Answer); i++ {
			if msg.Answer[i].Header().Rrtype == dns.TypeAAAA {
				// 删除 AAAA 记录
				msg.Answer = append(msg.Answer[:i], msg.Answer[i+1:]...)
				i--
			}
		}
	}

	res, _ := msg.Pack()
	// 将响应返回给客户端
	_, err = serverConn.WriteToUDP(res, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}
