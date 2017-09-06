package server

import (
	"github.com/astaxie/beego/logs"
	"github.com/gogo/protobuf/proto"
	"github.com/tonyzzp/ChaCha_Server/util"
	"github.com/tonyzzp/gocommon"
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
	userid      int
}

func NewClient(server *Server, conn *net.TCPConn) *Client {
	c := new(Client)
	c.server = server
	c.conn = conn
	c.alive = true
	c.init()
	c.sendingData = make(chan []byte, 100)
	return c
}

func (this *Client) init() {
	this.timeout(time.Minute)
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

func (this *Client) SendData(cmd int, msg proto.Message) {
	if this.alive {
		data := util.Cmd2Bytes(cmd, msg)
		this.sendingData <- data
	} else {
		logs.Warn("Client.SendData 用户已经无效")
		this.Close()
	}
}

func (this *Client) Close() {
	if !this.alive {
		return
	}
	this.alive = false
	this.server.onUserLogout(this)
	close(this.closeChan)
	this.conn.Close()
	close(this.sendingData)
	logs.Info("关闭客户端%v", this.conn.RemoteAddr())
	this = new(Client)
}

// 启动协程开始读取和发送数据
func (this *Client) Start() {
	go this.startRead()
	go this.startSend()
}

// 读取数据。  阻塞
func (this *Client) startRead() {
	for this.alive {
		size, e := util.ReadInt32(this.conn)
		if e != nil {
			logs.Info(e)
			this.Close()
			return
		}
		cmd, e := util.ReadInt32(this.conn)
		if e != nil {
			logs.Info(e)
			this.Close()
			return
		}
		content := make([]byte, size)
		_, e = io.ReadFull(this.conn, content)
		if e != nil {
			logs.Info(e)
			this.Close()
			return
		}
		this.timeout(time.Minute * 30)
		logs.Info("收到数据,%v cmd:%v size:%v", this.conn.RemoteAddr(), cmd, size)
		msg := new(Msg)
		msg.Client = this
		msg.Cmd = cmd
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
