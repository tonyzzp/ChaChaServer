package server

import (
	"github.com/astaxie/beego/logs"
	"github.com/gogo/protobuf/proto"
	"github.com/tonyzzp/ChaCha_Server/db"
	"github.com/tonyzzp/ChaCha_Server/protobeans"
	"time"
)

type Processor struct {
	server  *Server
	c       chan *Msg
	funcmap map[int]func(msg *Msg)
}

func NewProcessor(server *Server) *Processor {
	p := new(Processor)
	p.server = server
	p.c = make(chan *Msg, 100)
	p.funcmap = make(map[int]func(msg *Msg))

	p.funcmap[CMD_CLIENT_LOGIN] = p.onLogin
	p.funcmap[CMD_CLIENT_REGIST] = p.onRegist
	p.funcmap[CMD_CLIENT_SEND_TEXT_MSG] = p.onSendTextMessage

	return p
}

func (this *Processor) Start() {
	if this.c == nil {
		panic("需要调用NewProcessor来创建Processor")
	}
	go this.doStart()
}

//  开始阻塞处理数据
func (this *Processor) doStart() {
	for msg := range this.c {
		s := string(msg.Data)
		logs.Info("Processor: ip:%v cmd:%v len:%v", msg.Client.RemoteAdd(), msg.Cmd, len(s))
		f := this.funcmap[int(msg.Cmd)]
		if f == nil {
			logs.Warn("有消息无法处理 cmd:%v", msg.Cmd)
		} else {
			f(msg)
		}
	}
}

//  往队列里添加一条Msg
func (this *Processor) Append(msg *Msg) {
	this.c <- msg
}

func (this *Processor) Stop() {
	close(this.c)
}

////////////////////////////  命令处理函数

// 处理登录
func (this *Processor) onLogin(msg *Msg) {
	logs.Info("Processor.onLogin")
	if msg.Client.userid > 0 {
		logs.Warn("用户在登录状态下发送了登录命令，忽略")
		return
	}
	var bean protobeans.Login
	proto.Unmarshal(msg.Data, &bean)
	logs.Info(bean.String())
	u := db.FindUserByUserName(bean.Username)
	r := protobeans.LoginResponse{}
	if u != nil && u.PassWord == bean.Password {
		r.Ok = true
		msg.Client.userid = u.Id
		this.server.onUserLogin(msg.Client)
	} else {
		r.Ok = false
		r.Error = "用户名密码不正确"
	}
	logs.Info("回复 %v", r)
	msg.Client.SendData(CMD_CLIENT_LOGIN_RESP, &r)
}

// 注册
func (this *Processor) onRegist(msg *Msg) {
	logs.Info("Processor.onRegist")
	if msg.Client.userid > 0 {
		logs.Warn("用户在登录状态下发送了注册命令，忽略")
		return
	}
	var bean = protobeans.Regist{}
	proto.Unmarshal(msg.Data, &bean)
	logs.Info(bean)
	r := protobeans.RegistResponse{}
	if db.FindUserByUserName(bean.Username) != nil {
		r.Ok = false
		r.Error = "用户名已存在"
	} else {
		u := new(db.User)
		u.UserName = bean.Username
		u.PassWord = bean.Password
		ok := db.InsertUser(u)
		r.Ok = ok
		if ok {
			r.UserId = int32(u.Id)
			r.UserName = u.UserName
		}
	}
	logs.Info("回复 %v", r)
	msg.Client.SendData(CMD_CLIENT_REGIST_RESP, &r)
}

//  发文字消息
func (this *Processor) onSendTextMessage(msg *Msg) {
	logs.Info("Processor.onSendTextMessage")
	var bean = protobeans.TextMessage{}
	proto.Unmarshal(msg.Data, &bean)
	logs.Info(bean)
	ids := db.FindFriends(msg.Client.userid)
	receiver := int(bean.Receiver)
	found := false
	for _, id := range ids {
		if id == receiver {
			found = true
			break
		}
	}

	r := protobeans.TextMessgeResponse{}
	if !found {
		r.Error = "对方不是你的好友"
	} else {
		u := this.server.onlineUsers[receiver]
		if u == nil {
			r.Error = "对方不在线"
		} else {
			r.Ok = true

			m := protobeans.ReceiveTextMessage{}
			m.Sender = int32(msg.Client.userid)
			m.Content = bean.Content
			m.SendTime = time.Now().Unix()
			logs.Info("回复 %v", m)
			go u.SendData(CMD_SERVER_SEND_TEXT_MSG, &m)
		}
	}
	logs.Info("回复 %v", r)
	msg.Client.SendData(CMD_CLIENT_SEND_TEXT_MSG_RESP, &r)
}
