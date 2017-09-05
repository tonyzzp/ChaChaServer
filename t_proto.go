package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/tonyzzp/ChaCha_Server/protobeans"
)

func main() {
	t := protobeans.TextMessage{}
	t.Receiver = 123
	t.Content = "r u ok?"
	content, e := proto.Marshal(&t)
	fmt.Println(content)
	fmt.Println(e)
}
