package log

import (
	"log"
	"os"
)

type LoggerFactory struct {
	//log
	Logger *log.Logger
	//Error
	Error *log.Logger
}

func New() *LoggerFactory {
	return &LoggerFactory{
		Logger: nil,
		Error:  nil,
	}
}
func (logf *LoggerFactory) LogInit(logPath, errorPath string) error {
	//open file
	file, err := os.OpenFile(logPath+"log.log", os.O_APPEND|os.O_CREATE, 666)
	errFile, err := os.OpenFile(errorPath+"error.log", os.O_APPEND|os.O_CREATE, 666)
	if err != nil {
		return err
	}
	logf.Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	logf.Logger.Println("log init...")
	//set error
	logf.Error = log.New(errFile, "", log.LstdFlags|log.Lshortfile)
	logf.Error.Println("error log init...")
	return nil
}
