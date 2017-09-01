package server

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
)

type Server struct {
	listener *net.TCPListener
	running  bool
}

func NewServer(localAddress string) *Server {
	svr := new(Server)
	svr.running = true

	logs.GetLogger()

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
	client := NewClient(conn)
	client.Start()
}
