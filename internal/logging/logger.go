package logging

import (
	"log"
	"os"
)

type Logger struct {
	Warn  *log.Logger
	Info  *log.Logger
	Error *log.Logger
}

func NewLogger(name string) *Logger {
	fd, err := os.OpenFile(name+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return &Logger{
		Warn:  log.New(fd, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:  log.New(fd, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(fd, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
