package util

import (
	"github.com/tonyzzp/gocommon/bytesutil"
	"io"
)

//  从 r 中读取4字节的int。 高字节在前
func ReadInt32(r io.Reader) (int32, error) {
	header := make([]byte, 4)
	_, e := io.ReadFull(r, header)
	if e != nil {
		return 0, e
	}
	size := bytesutil.BytesToInt32(header)
	return size, nil
}
