package util

//wrap the channel
type SendRecvChan struct {
	send chan<- interface{}
	recv <-chan interface{}
}

func New() *SendRecvChan {
	channel := make(chan interface{})
	SRC := SendRecvChan{}
	SRC.send = channel
	SRC.recv = channel
	return &SRC
}

//send or recv must call by go routine.
func (ch *SendRecvChan) Send(i interface{}) {
	ch.send <- i
}
func (ch *SendRecvChan) Recv() interface{} {
	return <-ch.recv
}
