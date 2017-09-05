package main

import (
	"bufio"
	"fmt"
	"github.com/tonyzzp/gocommon/bytesutil"
	"net"
	"os"
	"strings"
)

func main() {
	ip, _ := net.ResolveTCPAddr("tcp", ":2626")
	conn, _ := net.DialTCP("tcp", nil, ip)
	r := bufio.NewReader(os.Stdin)
	for {
		s, _ := r.ReadString('\n')
		s = strings.TrimSpace(s)
		content := []byte(s)
		fmt.Println("content", content)
		data := bytesutil.Int32ToBytes(int32(len(content)))
		fmt.Println("header", data)
		data = append(data, content...)
		fmt.Println(data)
		conn.Write(data)
	}
}
