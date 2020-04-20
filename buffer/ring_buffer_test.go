package buffer

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	buf := New(2097152)
	buf.Write([]byte("12345678123456781234567812345678123456781234567812345678123456781234567812345678"))
	bytes, _ := buf.ReadView(12)
	fmt.Println(bytes)
	//bytes[0] = 64
	//fmt.Println(bytes)
	fmt.Println(buf.size)
	b := make([]byte, 8)
	buf.Read(b)
	//fmt.Println(string(b))
	fmt.Println(buf.size)
	fmt.Println(buf.size)
}
