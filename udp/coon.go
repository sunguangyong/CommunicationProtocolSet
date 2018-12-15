package udp

import (
	"encoding/json"
	"fubangyun.com/basearch/goutil/logger"
	"fubangyun.com/basearch/manna-stream/device/protocol"
	"net"
	"sync"
	"time"
)

type ConnHandler interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}


const (
	DeadLineSecs      = 3600
	defaultBufferSize = 1024
	timeout           = 60 * 5
)

type Messager struct {
	appType string
	dtuType string
	dtuId   int64

	lastTime  string
	startTime string

	rcvNum int
	pktNum int
	nsqNum int
	ackNum int

	mutex *sync.Mutex

	conn ConnHandler

	parser interface{} //parser.Parser
	proto  protocol.Protocoler

	messages chan []byte

	remote string
	active int

}


type Conn struct {
	appType string //server app name
	conn    ConnHandler
	inBuf   []byte

	msger *Messager

	isRecognized bool
	flag         chan byte
}

func newConn(conn ConnHandler, name string) *Conn {
	c := &Conn{
		appType:      name,
		conn:         conn,
		inBuf:        make([]byte, defaultBufferSize),
		msger:        nil,
		isRecognized: false,
		flag:         make(chan byte),
	}
	return c
}

func (c *Conn) Start() {
	defer func() {
		if err := recover(); err != nil {
			logger.Println(c.appType+" Recover(Conn.Start):", err)
		}

		logger.Println(c.appType + " Disconnected :" + c.RemoteAddr())
		c.conn.Close()
	}()


	go c.heartBeating(timeout)
	c.readPump()
}

func (c *Conn) LocalAddr() string {
	return c.conn.LocalAddr().String()
}

func (c *Conn) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

type DevOnline struct {
	ClientId   int64  `json:"client_id"`
	ClientType string `json:"client_type"`
	NetPro     string `json:"net_pro"`
	OnlineTime string `json:"online_time"`
}

func (dev *DevOnline) ToJson() string {
	data, err := json.Marshal(dev)
	if err != nil {
		return ""
	}
	return string(data)
}

func (c *Conn) heartBeating(timeout int) {
	for {
		select {
		case <-c.flag:
			break
		case <-time.After(time.Second * time.Duration(timeout)):
			c.conn.Close()
		}
	}
}

func (c *Conn) readPump() {
	ip := c.RemoteAddr()
	n, err := c.conn.Read(c.inBuf)
	if err != nil {
		logger.Println(ip, " Read Fist Data error: ", err)
		return
	}
	logger.Println(ip, "Receive First data:", c.inBuf[:n], " app_type: ", c.appType)

}
