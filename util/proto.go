package util

import (
	"github.com/gogo/protobuf/proto"
	"github.com/tonyzzp/gocommon/bytesutil"
)

// 将cmd组装成最终要发送的包的格式：len,cmd,data
func Cmd2Bytes(cmd int, msg proto.Message) []byte {
	d, _ := proto.Marshal(msg)
	l := len(d)
	data := append(bytesutil.Int32ToBytes(int32(l)), bytesutil.Int32ToBytes(int32(cmd))...)
	data = append(data, d...)
	return data
}
