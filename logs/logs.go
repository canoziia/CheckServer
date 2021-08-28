package logs

import (
	"fmt"
	"log"
	"os"
)

var (
	GlobalLogger *Logger
)

type Info struct {
	Code    int
	Message string
}

type Logger struct {
	*log.Logger
	File *os.File
}

func InitGlobalLogger(filePath string) {
	GlobalLogger = NewLogger(filePath)
}

func NewLogger(filePath string) (logger *Logger) {
	logger = new(Logger)
	var err error
	logger.File, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(fmt.Sprintf("Load LogFile %s Failed: %s", filePath, err.Error()))
	}
	logger.Logger = log.New(logger.File, "[CheckServer] ", log.LstdFlags|log.LUTC) // |log.Lshortfile
	return
}

func (lgr *Logger) Record(info Info) {
	statusList := []string{"[failed]", "[info]", "[error]"}
	log.Print(fmt.Sprintf("%s %s", statusList[info.Code], info.Message))
	lgr.Print(fmt.Sprintf("%s %s", statusList[info.Code], info.Message))
}

func (lgr *Logger) AllPrint(text string) {
	lgr.Print(text)
	log.Print(text)
}
