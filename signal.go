package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Signal() {

	ch := make(chan os.Signal, 1)

	s := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	}

	signal.Notify(ch, s...)

	sig := <-ch

	fmt.Println(sig)

}
