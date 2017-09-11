package main

import (
	"fmt"
	"github.com/tonyzzp/gocommon/bytesutil"
)

func main() {
	i := 129
	content := bytesutil.Int32ToBytes(int32(i))
	fmt.Println(content)
}
