package main

import (
	"grtp/grtp"
	"testing"
	"time"
)

/**
>go test -v -bench="." -benchmem ./main_test.go ./main.go
goos: windows
goarch: amd64
BenchmarkTestSync-4          100          12047401 ns/op             116 B/op          0 allocs/op
BenchmarkTestAsync-4         591          10753934 ns/op               9 B/op          0 allocs/op

*/
func BenchmarkTestSync(b *testing.B) {
	p, _ := grtp.New(8)
	go p.Run()
	task := grtp.NewTask(func() {
		time.Sleep(time.Millisecond * 100)
		//fmt.Println("task1.")
	})

	for i := 0; i < b.N; i++ {
		p.Commit(task)
	}

}
func BenchmarkTestAsync(b *testing.B) {
	p, _ := grtp.NewAsync(8, 40)
	go p.Run()
	task := grtp.NewTask(func() {
		time.Sleep(time.Millisecond * 100)
		//fmt.Println("task1.")
	})

	for i := 0; i < b.N; i++ {
		p.Commit(task)
	}
}
