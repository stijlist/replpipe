package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()
	repl, err := net.Dial("tcp", os.Args[1])
	defer repl.Close()
	if err != nil {
		log.Panic(fmt.Errorf("failed to connect to repl: %w", err))
	}
	if err := unix.Mkfifo(".repl-pipe", 0666); err != nil {
		log.Panic(fmt.Errorf("mkfifo: %w", err))
	}
	defer func() {
		if err := os.Remove(".repl-pipe"); err != nil {
			log.Println(err)
		}
	}()
	f, err := os.OpenFile(".repl-pipe", os.O_RDWR, os.ModeNamedPipe)
	defer f.Close()
	if err != nil {
		log.Panic(err)
	}
	fifo := bufio.NewReader(f)
	replAndStdout := io.MultiWriter(repl, os.Stdout)
	go func() {
		// this goroutine should exit when `repl` is closed
		if _, err := io.Copy(os.Stdout, repl); err != nil {
			log.Println(err)
			return
		}
	}()
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			for {
				n, err := fifo.WriteTo(replAndStdout)
				if err != nil && !errors.Is(err, io.EOF) {
					log.Panic(fmt.Errorf("failed to write: %w", err))
				}
				if n == 0 {
					log.Printf("%s: zero sized write", repl)
				}
			}
		}
	}()
	<-ctx.Done()
}
