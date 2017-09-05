package server

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/tonyzzp/gocommon"
	"github.com/tonyzzp/gocommon/bytesutil"
	"io"
	"net"
	"time"
)

type Client struct {
	server      *Server
	conn        *net.TCPConn
	alive       bool
	closeChan   chan int
	sendingData chan []byte
}

func NewClient(server *Server, conn *net.TCPConn) *Client {
	c := new(Client)
	c.server = server
	c.conn = conn
	c.alive = true
	c.init()
	c.sendingData = make(chan []byte)
	return c
}

func (this *Client) init() {
	this.timeout(10 * time.Second)
}

func (this *Client) RemoteAdd() net.Addr {
	return this.conn.RemoteAddr()
}

// 设置一个定时器， d时间之后关闭client
func (this *Client) timeout(d time.Duration) {
	if this.closeChan != nil {
		close(this.closeChan)
	}
	this.closeChan = gocommon.TimeOut(d, func() {
		logs.Info("客户端超时 %v", this.conn.RemoteAddr())
		this.Close()
	})
}

func (this *Client) SendData(data []byte) {
	this.sendingData <- data
}

func (this *Client) Close() {
	if !this.alive {
		return
	}
	this.alive = false
	close(this.closeChan)
	this.conn.Close()
	close(this.sendingData)
	logs.Info("关闭客户端%v", this.conn.RemoteAddr())
}

// 启动协程开始读取和发送数据
func (this *Client) Start() {
	go this.startRead()
	go this.startSend()
}

// 读取数据。  阻塞
func (this *Client) startRead() {
	for this.alive {
		header := make([]byte, 4)
		_, e := io.ReadFull(this.conn, header)
		if e != nil {
			logs.Info(e)
			this.Close()
			return
		}
		size := bytesutil.BytesToInt32(header)
		fmt.Println(size)
		content := make([]byte, size)
		_, e = io.ReadFull(this.conn, content)
		if e != nil {
			logs.Info(e)
			this.Close()
			return
		}
		this.timeout(time.Minute)
		logs.Info("收到数据,%v size:%v", this.conn.RemoteAddr(), size)
		msg := new(Msg)
		msg.Client = this
		msg.Data = content
		this.server.onClientMsg(msg)
	}
}

// 从缓冲队列里取数据出来发送出去。 阻塞
func (this *Client) startSend() {
	for data := range this.sendingData {
		this.conn.Write(data)
	}
}
