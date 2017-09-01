package main

import (
	"fmt"
	"net"
	"testing"
)

func Test_client(t *testing.T) {
	addr, _ := net.ResolveTCPAddr("tcp", ":2626")
	conn, e := net.DialTCP("tcp", nil, addr)
	fmt.Println(e)
	conn.Write([]byte("章治鹏"))
	conn.Write([]byte("zzp"))
	conn.Close()
}

func Test_readint(t *testing.T) {
	var i int32 = 4000
	b := []byte{
		byte((i >> 24) & 0xff),
		byte((i >> 16) & 0xff),
		byte((i >> 8) & 0xff),
		byte(i & 0xff),
	}
	fmt.Println(i)
	fmt.Println(b)

	fmt.Println(int32(b[0]<<24), int32(b[1]<<16), b[2]<<8, int32(b[3]))

}
