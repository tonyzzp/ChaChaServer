package server

import (
	"github.com/astaxie/beego/logs"
	"github.com/gogo/protobuf/proto"
	"github.com/tonyzzp/ChaCha_Server/util"
	"github.com/tonyzzp/gocommon"
	"github.com/tonyzzp/gocommon/bytesutil"
	"io"
	"net"
	"sync"
	"time"
)

var (
	mux = sync.Mutex{}
	id  = 0
)

func sequenceId() int {
	v := 0
	mux.Lock()
	id++
	v = id
	mux.Unlock()
	return v
}

type Client struct {
	server      *Server
	conn        *net.TCPConn
	alive       bool
	closeChan   chan int
	sendingData chan []byte
	userid      int
	sequence    int
}

func NewClient(server *Server, conn *net.TCPConn) *Client {
	c := new(Client)
	c.sequence = sequenceId()
	c.server = server
	c.conn = conn
	c.alive = true
	c.sendingData = make(chan []byte, 100)

	c.timeout(time.Minute)
	return c
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
		if msg != nil {
			data := util.Cmd2Bytes(cmd, msg)
			this.sendingData <- data
		} else {
			data := bytesutil.Int32ToBytes(0)
			data = append(data, bytesutil.Int32ToBytes(int32(cmd))...)
			this.sendingData <- data
		}
	} else {
		logs.Warn("Client.SendData 用户已经无效")
		this.Close()
	}
}

// 将client标记为已关闭。  并不会真的关闭connection,而是等待所有队列里的数据写完后才会真的关闭connection
func (this *Client) Close() {
	if !this.alive {
		return
	}
	this.alive = false
	close(this.closeChan)
	close(this.sendingData)
	this.server.onUserLogout(this)
	logs.Info("即将关闭客户端%v", this.conn.RemoteAddr())
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
	this.conn.Close()
	logs.Info("%v 数据全部发送完成，close connection", this.conn.RemoteAddr())
}
