package main

import (
	"fmt"
	"github.com/tonyzzp/gocommon/bytesutil"
	"testing"
)

func Test_readint(t *testing.T) {
	var i int32 = 2147483647
	b := []byte{
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
	fmt.Println(i)
	fmt.Println(b)

	i = int32(int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3]))
	fmt.Println(i)

	b = bytesutil.Int32ToBytes(10000123)
	fmt.Println(b)
	fmt.Println(bytesutil.BytesToInt32(b))
}
