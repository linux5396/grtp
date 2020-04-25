package log

import (
	"testing"
	"time"
)

func TestLoggerFactory_LogInit(t *testing.T) {
	LoggerFactory := New()
	_ = LoggerFactory.LogInit("./log/", "./log/")
	LoggerFactory.Logger.Println("main.log")
	LoggerFactory.Error.Println("occurred fa")
	LoggerFactory.Error.Println("main.error.log")
	for i := 0; i < 10; i++ {
		if i == 4 {
			LoggerFactory.Logger.Println("gc...")
			LoggerFactory.LogDestroy() //test destroyed
		}
		time.Sleep(time.Second)
		LoggerFactory.Logger.Println("main.log")
		LoggerFactory.Error.Println("occurred fa")
	}
}
