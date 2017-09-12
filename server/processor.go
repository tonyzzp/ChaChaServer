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
	p.funcmap[CMD_CLIENT_QUERY_USER] = p.onSearchUser
	p.funcmap[CMD_CLIENT_ADD_FRIEND] = p.onAddFriend
	p.funcmap[CMD_CLIENT_ADD_FRIEND_CONFIRM] = p.onAddFriendConfirm
	p.funcmap[CMD_CLIENT_FRIENDS_LIST] = p.onFriendsList
	p.funcmap[CMD_CLIENT_REMOVE_FRIEND] = p.onRemoveFriend

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
		exist := this.server.onlineUsers[u.Id]
		if exist != nil {
			exist.SendData(CMD_SERVER_KICKOUT, nil)
			exist.Close()
		}
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
	r.Sequence = bean.Sequence
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
			u.SendData(CMD_SERVER_SEND_TEXT_MSG, &m)
		}
	}
	logs.Info("回复 %v", r)
	msg.Client.SendData(CMD_CLIENT_SEND_TEXT_MSG_RESP, &r)
}

// 查找用户
func (this *Processor) onSearchUser(msg *Msg) {
	logs.Info("Processor.onSearchUser")
	bean := new(protobeans.SearchUser)
	proto.Unmarshal(msg.Data, bean)
	logs.Info(bean)
	u := db.FindUserByUserName(bean.UserName)
	r := new(protobeans.SearchUserResp)
	r.UserName = bean.UserName
	if u != nil {
		r.UserId = int32(u.Id)
		r.Sex = int32(u.Sex)
		r.Head = u.Head
	}
	msg.Client.SendData(CMD_CLIENT_QUERY_USER_RESP, r)
}

//请求添加好友
func (this *Processor) onAddFriend(msg *Msg) {
	logs.Info("Processor.onAddFriend")
	bean := new(protobeans.AddFriend)
	proto.Unmarshal(msg.Data, bean)
	r := new(protobeans.AddFriendResp)
	friend := this.server.onlineUsers[int(bean.FriendId)]
	if friend != nil {
		r.Ok = true

		r2 := new(protobeans.AddFriend)
		r2.FriendId = int32(msg.Client.userid)
		r2.Message = bean.Message
		friend.SendData(CMD_SERVER_ADD_FRIEND, r2)
	} else {
		r.Ok = false
		r.Error = "对方不在线"
	}
	msg.Client.SendData(CMD_CLIENT_ADD_FRIEND_RESP, r)
}

//通过或拒绝添加好友
func (this *Processor) onAddFriendConfirm(msg *Msg) {
	logs.Info("Processor.onADdFriendConfirm")
	bean := new(protobeans.AddFriendConfirm)
	proto.Unmarshal(msg.Data, bean)
	logs.Info(bean)
	if bean.Ok {
		friend := db.FindUserById(int(bean.FriendId))
		if friend != nil {
			db.AddFriend(msg.Client.userid, friend.Id)
			db.AddFriend(friend.Id, msg.Client.userid)

			r := new(protobeans.NewFriend)
			r.FriendId = int32(friend.Id)
			r.UserName = friend.UserName
			r.Sex = int32(friend.Sex)
			r.Head = friend.Head
			msg.Client.SendData(CMD_SERVER_NEW_FRIEND, r)

			fc := this.server.onlineUsers[int(bean.FriendId)]
			if fc != nil {
				u := db.FindUserById(msg.Client.userid)
				r := new(protobeans.NewFriend)
				r.FriendId = int32(u.Id)
				r.UserName = u.UserName
				r.Sex = int32(u.Sex)
				r.Head = u.Head
				fc.SendData(CMD_SERVER_NEW_FRIEND, r)
			}
		}
	}
}

// 请求好友列表
func (this *Processor) onFriendsList(msg *Msg) {
	logs.Info("Processor.onFriendsList")
	ids := db.FindFriends(msg.Client.userid)
	logs.Info(ids)
	users := db.FetchUsrs(ids)
	logs.Info(users)
	r := new(protobeans.FriendsList)
	r.Friends = make([]*protobeans.Friend, len(ids))
	for i := 0; i < len(users); i++ {
		f := new(protobeans.Friend)
		r.Friends[i] = f
		u := users[i]
		f.UserId = int32(u.Id)
		f.UserName = u.UserName
		f.Sex = int32(u.Sex)
		f.Head = u.Head
	}
	msg.Client.SendData(CMD_CLIENT_FRIENDS_LIST_RESP, r)
}

// 删除好友
func (this *Processor) onRemoveFriend(msg *Msg) {
	logs.Info("Processor.onRemoveFriend")
	r := new(protobeans.RemoveFriend)
	proto.Unmarshal(msg.Data, r)
	db.RemoveFriend(msg.Client.userid, int(r.UserId))
	msg.Client.SendData(CMD_CLIENT_REMOVE_FRIEND_RESP, r)
}
