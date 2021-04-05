package logging

import (
	"io"
	"log"
)

const (
	LdefaultFlags = log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix
	LstdOutPrefix = "[INFO] "
	LstdErrPrefix = "[ERROR] "
)

func NewStdOutLogger(out io.Writer) *log.Logger {
	return log.New(out, LstdOutPrefix, LdefaultFlags)
}

func NewStdErrLogger(out io.Writer) *log.Logger {
	return log.New(out, LstdErrPrefix, LdefaultFlags)
}
