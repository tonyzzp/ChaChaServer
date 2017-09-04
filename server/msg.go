package server

//  收到客户端消息后，将以Msg的形式封装，然后发送给Processor进行处理
type Msg struct {
	Client *Client
	Data   []byte
}
