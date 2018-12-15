package udp

import (
	_ "fmt"
	"fubangyun.com/basearch/goutil/logger"
	"net"
)

type UdpConn struct {
	addr            *net.UDPAddr
	inBufs          chan []byte
	UdpServerSocket *net.UDPConn
}

func (c *UdpConn) Read(data []byte) (int, error) {
	tmpBuf := <-c.inBufs
	logger.Printf("udp read!!!!! data:%v", tmpBuf)
	copy(data, tmpBuf)
	return len(tmpBuf), nil
}

func (c *UdpConn) Write(data []byte) (int, error) {
	logger.Printf("udp wite to addr:%v, data:%v", c.addr, data)
	return c.UdpServerSocket.WriteToUDP(data, c.addr)
}

func (c *UdpConn) Close() error {
	return nil
}

func (c *UdpConn) LocalAddr() net.Addr {
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:1911")
	return addr
}

func (c *UdpConn) RemoteAddr() net.Addr {
	return c.addr
}
