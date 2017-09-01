package server

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/tonyzzp/gocommon"
	"net"
	"time"
)

type Client struct {
	conn      *net.TCPConn
	alive     bool
	closeChan chan int
}

func NewClient(conn *net.TCPConn) *Client {
	c := new(Client)
	c.conn = conn
	c.alive = true
	c.init()
	return c
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

func (this *Client) init() {
	this.timeout(10 * time.Second)
}

func (this *Client) Close() {
	close(this.closeChan)
	this.conn.Close()
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
		content := make([]byte, 100)
		count, e := this.conn.Read(content)
		if e != nil {
			fmt.Println(e)
			this.Close()
			return
		}
		this.timeout(time.Minute)
		logs.Info("收到数据 %v", string(content[:count]))
	}
}

// 从缓冲队列里取数据出来发送出去。 阻塞
func (this *Client) startSend() {

}
