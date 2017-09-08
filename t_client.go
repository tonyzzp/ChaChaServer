package main

import (
	"bufio"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/tonyzzp/ChaCha_Server/protobeans"
	"github.com/tonyzzp/ChaCha_Server/server"
	"github.com/tonyzzp/ChaCha_Server/util"
	"github.com/tonyzzp/gocommon/bytesutil"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var cmds = make(map[string]func(s string))
var conn *net.TCPConn

func init() {
	cmds["reg"] = reg
	cmds["login"] = login
	cmds["send"] = send
	cmds["search"] = search
	cmds["add"] = add
	cmds["confirm"] = confirm
	cmds["friends"] = friends
	cmds["remove"] = remove
}

func reg(s string) {
	strs := strings.Split(s, " ")
	userName := strs[0]
	passWord := strs[1]
	r := protobeans.Regist{}
	r.Username = userName
	r.Password = passWord
	data := util.Cmd2Bytes(server.CMD_CLIENT_REGIST, &r)
	conn.Write(data)
}

func login(s string) {
	strs := strings.Split(s, " ")
	userName := strs[0]
	passWord := strs[1]
	r := new(protobeans.Login)
	r.Username = userName
	r.Password = passWord
	data := util.Cmd2Bytes(server.CMD_CLIENT_LOGIN, r)
	conn.Write(data)
}

func send(s string) {
	index := strings.Index(s, " ")
	id, _ := strconv.Atoi(s[:index])
	msg := s[index+1:]
	r := new(protobeans.TextMessage)
	r.Receiver = int32(id)
	r.Content = msg
	data := util.Cmd2Bytes(server.CMD_CLIENT_SEND_TEXT_MSG, r)
	conn.Write(data)
}

func search(s string) {
	r := new(protobeans.SearchUser)
	r.UserName = s
	data := util.Cmd2Bytes(server.CMD_CLIENT_QUERY_USER, r)
	conn.Write(data)
}

func add(s string) {
	r := new(protobeans.AddFriend)
	i, _ := strconv.Atoi(s)
	r.FriendId = int32(i)
	r.Message = "我想加你啊" + s
	data := util.Cmd2Bytes(server.CMD_CLIENT_ADD_FRIEND, r)
	conn.Write(data)
}

func confirm(s string) {
	r := new(protobeans.AddFriendConfirm)
	r.Ok = true
	i, _ := strconv.Atoi(s)
	r.FriendId = int32(i)
	data := util.Cmd2Bytes(server.CMD_CLIENT_ADD_FRIEND_CONFIRM, r)
	conn.Write(data)
}

func friends(s string) {
	data := append(bytesutil.Int32ToBytes(int32(0)), bytesutil.Int32ToBytes(int32(server.CMD_CLIENT_FRIENDS_LIST))...)
	conn.Write(data)
}

func remove(s string) {
	r := new(protobeans.RemoveFriend)
	i, _ := strconv.Atoi(s)
	r.UserId = int32(i)
	data := util.Cmd2Bytes(server.CMD_CLIENT_REMOVE_FRIEND, r)
	conn.Write(data)
}

func read() {
	for {
		l, e := util.ReadInt32(conn)
		if e != nil {
			panic(e)
		}
		cmd, _ := util.ReadInt32(conn)
		buf := make([]byte, l)
		io.ReadFull(conn, buf)
		switch int(cmd) {
		case server.CMD_CLIENT_LOGIN_RESP:
			b := new(protobeans.LoginResponse)
			proto.Unmarshal(buf, b)
			fmt.Println("登录回应:", b.Ok, b.Error)
		case server.CMD_CLIENT_REGIST_RESP:
			b := new(protobeans.RegistResponse)
			proto.Unmarshal(buf, b)
			fmt.Println("注册回应:", b.Ok, b.Error, b.UserId)
		case server.CMD_CLIENT_SEND_TEXT_MSG_RESP:
			b := new(protobeans.TextMessgeResponse)
			proto.Unmarshal(buf, b)
			fmt.Println("发消息回应:", b.Ok, b.Error)
		case server.CMD_SERVER_SEND_TEXT_MSG:
			b := new(protobeans.ReceiveTextMessage)
			proto.Unmarshal(buf, b)
			fmt.Println("收到消息:", b.Sender, b.Content, time.Unix(b.SendTime, 0).Format("2006-01-02 15:04:05"))
		case server.CMD_SERVER_KICKOUT:
			fmt.Println("异地登录")
		case server.CMD_CLIENT_QUERY_USER_RESP:
			b := new(protobeans.SearchUserResp)
			proto.Unmarshal(buf, b)
			fmt.Println("你查询的用户", b)
		case server.CMD_CLIENT_ADD_FRIEND_RESP:
			b := new(protobeans.AddFriendResp)
			proto.Unmarshal(buf, b)
			fmt.Println("添加用户返回", b.Ok, b.Error)
		case server.CMD_SERVER_ADD_FRIEND:
			b := new(protobeans.AddFriend)
			proto.Unmarshal(buf, b)
			fmt.Println("有人想加你为好友：", b.FriendId, b.Message)
		case server.CMD_SERVER_NEW_FRIEND:
			b := new(protobeans.NewFriend)
			proto.Unmarshal(buf, b)
			fmt.Println("新好友", b.FriendId, b.UserName)
		case server.CMD_CLIENT_FRIENDS_LIST_RESP:
			b := new(protobeans.FriendsList)
			proto.Unmarshal(buf, b)
			fmt.Println("好友列表:")
			for _, v := range b.Friends {
				fmt.Println(v.UserId, v.UserName)
			}
		case server.CMD_CLIENT_REMOVE_FRIEND_RESP:
			b := new(protobeans.RemoveFriend)
			proto.Unmarshal(buf, b)
			fmt.Println("删除好友成功:", b.UserId)
		default:
			fmt.Println("命令未处理:", strconv.FormatInt(int64(cmd), 16))
		}
	}
}

func main() {
	ip, _ := net.ResolveTCPAddr("tcp", ":2626")
	fmt.Println("开始连接", ip)
	conn, _ = net.DialTCP("tcp", nil, ip)
	if conn != nil {
		fmt.Println("连接服务器成功")
		fmt.Println(`1.注册 reg:NAME PASSWORD
2.登录 login:NAME PASSWORD
3.发消息 send:RECEIVER hello!
4.退出 exit
5.查询用户 search:NAME
5.添加好友 add:ID
6.通过请求 confirm:ID
7.好友列表 friends:
8.删除好友 remove:ID
`)
	} else {
		fmt.Println("连接失败")
		return
	}
	go read()
	r := bufio.NewReader(os.Stdin)
	for {
		s, _ := r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "exit" {
			conn.Close()
			break
		}
		strs := strings.Split(s, ":")
		f := cmds[strs[0]]
		if f != nil {
			f(strs[1])
		} else {
			fmt.Println("命令错误，正确命令如：   login:aaa 111")
		}
	}
}
