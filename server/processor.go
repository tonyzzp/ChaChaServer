package server

import (
	"github.com/astaxie/beego/logs"
)

type Processor struct {
	c chan *Msg
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.c = make(chan *Msg)
	return p
}

func (this *Processor) Start() {
	if this.c == nil {
		panic("需要调用NewProcessor来创建Processor")
	}
	go this.doStart()
}

func (this *Processor) doStart() {
	for msg := range this.c {
		s := string(msg.Data)
		a := msg.Client.RemoteAdd()
		logs.Info("Processor: %v : %v", a, s)
		if s == "c" {
			msg.Client.Close()
		}
	}
}

//  往队列里添加一条Msg
func (this *Processor) Append(msg *Msg) {
	this.c <- msg
}

func (this *Processor) Stop() {
	close(this.c)
}
