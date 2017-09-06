package main

import (
	"bufio"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/tonyzzp/ChaCha_Server/protobeans"
	"github.com/tonyzzp/ChaCha_Server/server"
	"github.com/tonyzzp/ChaCha_Server/util"
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

func read() {
	for {
		l, _ := util.ReadInt32(conn)
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
4.退出 exit`)
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
