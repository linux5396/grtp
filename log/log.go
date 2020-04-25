package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"unsafe"
)

type LoggerFactory struct {
	//log
	Logger *log.Logger
	//Error
	Error *log.Logger
	//register hook
	cancel context.CancelFunc
}

func New() *LoggerFactory {
	return &LoggerFactory{
		Logger: nil,
		Error:  nil,
	}
}

/**
 * path should be exists.
 * 0666 -> owner group root r=4 w=2 x=1
 * 6->rw
 */
func (logf *LoggerFactory) LogInit(logPath, errorPath string) error {
	//open file
	file, err := os.OpenFile(logPath+"log.log", os.O_APPEND|os.O_CREATE, 0666)
	errFile, err := os.OpenFile(errorPath+"error.log", os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background()) //register context
	fmt.Printf("log:%p", unsafe.Pointer(&ctx))
	logf.cancel = cancel //make cancel
	go func() {          //register hook
		select {
		case <-ctx.Done():
			file.Close()
			errFile.Close()
		}
	}()
	logf.Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	//chose console to print
	log.Println("log init...")
	//set error
	logf.Error = log.New(errFile, "", log.LstdFlags|log.Lshortfile)
	//chose console to print
	log.Println("error log init...")
	return nil
}

//using the context to close the file .
func (logf *LoggerFactory) LogDestroy() {
	logf.cancel() //cancel
}
