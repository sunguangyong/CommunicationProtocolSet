package main

import (
	_ "fmt"
	"fubangyun.com/basearch/CommunicationProtocolSet/udp"
)

func main() {
	udp_server := &udp.UdpServer{"0.0.0.0:5000", "NBDQHZ"} //NB模块电气火灾
	udp_server.Run()
}