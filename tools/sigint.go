package tools

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func HandleSignalInterrupt() {
	log.Println("Please press ctrl+c to exit.")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()
}
