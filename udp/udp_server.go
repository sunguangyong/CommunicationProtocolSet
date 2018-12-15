package udp

import (
	_ "fmt"
	"fubangyun.com/basearch/goutil/logger"
	"fubangyun.com/basearch/manna-stream/common"
	"net"
)

var UdpConnPool map[string]*UdpConn

func init() {
	UdpConnPool = make(map[string]*UdpConn)
}

type UdpServer struct {
	Host string
	Name string
}

func (svr *UdpServer) Run() {
	exit := make(chan bool, 1)

	addr, _ := net.ResolveUDPAddr("udp4", svr.Host)
	var err error

	var UdpServerSocket *net.UDPConn
	UdpServerSocket, err = net.ListenUDP("udp4", addr)
	if err != nil {
		logger.Println("Udp server listen failed!", err)
		return
	}
	defer UdpServerSocket.Close()

	go func() {
		for {
			logger.Println("<<<<<<<<<<<<<<<<<")
			data := make([]byte, 4096)
			n, remoteAddr, err := UdpServerSocket.ReadFromUDP(data)

			if err != nil {
				logger.Println("Udp Server read data error!", err)
				continue
			}

			go svr.ReceiveData(UdpServerSocket, remoteAddr, data[:n])

			logger.Println(">>>>>>>>>>>>>>>>>")
		}
	}()

	logger.Println(common.NowTime() + " Tcp Server Scucess Start on: " + svr.Host)
	<-exit
	logger.Println(common.NowTime() + " Tcp Server Died on: " + svr.Host)
}

func (svr *UdpServer) ReceiveData(UdpServerSocket *net.UDPConn, remoteAddr *net.UDPAddr, data []byte) {
	logger.Printf("Accept UDP Connection, Receive Data, remoteAddr=%s, len=%d, data(hex)=[% x]\n", remoteAddr, len(data), data)
	addr := remoteAddr.String()
	if _, ok := UdpConnPool[addr]; !ok {
		UdpConnPool[addr] = &UdpConn{
			addr:            remoteAddr,
			inBufs:          make(chan []byte),
			UdpServerSocket: UdpServerSocket,
		}
		c := newConn(UdpConnPool[addr], svr.Name)
		c.Start()
	}
	UdpConnPool[addr].inBufs <- data
}
