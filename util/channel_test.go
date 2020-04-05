package util

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ch := New()
	go ch.Send("some data")
	go fmt.Println(ch.Recv())

	blocked := make(chan struct{})
	go func() {
		fmt.Println("enter blocked")
		time.Sleep(time.Millisecond * 2000)
		blocked <- struct{}{} //release
		fmt.Println("exit blocked")
	}()
	<-blocked
	time.Sleep(time.Millisecond)
	fmt.Println("main run")
}
