package server

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
)

type Server struct {
	listener    *net.TCPListener
	running     bool
	processor   *Processor
	onlineUsers map[int]*Client
}

func NewServer(localAddress string) *Server {
	svr := new(Server)
	svr.running = true
	svr.processor = NewProcessor(svr)
	svr.onlineUsers = make(map[int]*Client)

	ip, e := net.ResolveTCPAddr("tcp", localAddress)
	if e != nil {
		panic(fmt.Sprint("获取ip失败", e))
	}

	svr.listener, e = net.ListenTCP("tcp", ip)
	if e != nil {
		panic(fmt.Sprint("启动服务失败", e))
	}
	logs.Info("启动服务成功 ip: %v", ip)

	return svr
}

// 启动监听。  方法会阻塞
func (this *Server) Start() {
	if this.listener == nil {
		panic("请使用NewServer来进行创建")
	}
	this.processor.Start()
	for this.running {
		conn, e := this.listener.AcceptTCP()
		if e != nil {
			logs.Warn("客户端接入失败 %v", e)
		}
		logs.Info("客户端接入 %v", conn.RemoteAddr())
		this.onAccept(conn)
	}
}

// 收到一个客户端链接
func (this *Server) onAccept(conn *net.TCPConn) {
	client := NewClient(this, conn)
	client.Start()
}

// 有客户端发来一条消息。调用此方法会将msg放到队列，等待处理
func (this *Server) onClientMsg(msg *Msg) {
	this.processor.Append(msg)
}

func (this *Server) onUserLogin(client *Client) {
	this.onlineUsers[client.userid] = client
}

func (this *Server) onUserLogout(client *Client) {
	delete(this.onlineUsers, client.userid)
}

func (this *Server) isUserOnline(userid int) bool {
	return this.onlineUsers[userid] != nil
}
