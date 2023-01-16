package main

import (
	"github.com/MoeQing-Network/MoeDNS/config"
	"github.com/MoeQing-Network/MoeDNS/dns"
)

func main() {
	config.Init()
	dns.Start()
}
