package tools

import (
	"log"
	"os"
)

func CreateDirIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0655); err != nil {
			log.Fatalf("could not create static file path %q: %v\n", path, err)
		}
	}
}

func CreateOrOpenFile(name string) *os.File {
	fd, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return fd
}
