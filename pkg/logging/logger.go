package logging

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	Warn  *log.Logger
	Info  *log.Logger
	Error *log.Logger
}

func NewStdoutLogger() *Logger {
	return &Logger{
		Warn:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func NewLogFile(name string) *os.File {
	fd, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return fd
}

func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = NewLogFile("webserver.log")
	}
	return &Logger{
		Warn:  log.New(w, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:  log.New(w, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(w, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
