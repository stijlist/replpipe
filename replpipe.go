package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	addr := os.Args[1]
	repl, err := net.Dial("tcp", addr)
	defer repl.Close()
	if err != nil {
		log.Fatalf("failed to connect to repl: %s", err)
	}
	if err := unix.Mkfifo(".repl-pipe", 0666); err != nil {
		log.Fatalf("mkfifo: %s", err)
	}
	defer func() {
		if err := os.Remove(".repl-pipe"); err != nil {
			log.Println(err)
		}
	}()
	f, err := os.OpenFile(".repl-pipe", os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(f)
	fifo := bufio.NewReader(f)
	// print input to stdout and send to repl
	replAndStdout := io.MultiWriter(repl, os.Stdout)
	go func() {
		if _, err := io.Copy(os.Stdout, repl); err != nil {
			log.Fatal(err)
		}
	}()
	for {
		n, err := fifo.WriteTo(replAndStdout)
		if err != nil && !errors.Is(err, io.EOF) {
			log.Fatalf("failed to write: %s", err)
		}
		if n == 0 {
			log.Printf("zero sized write to %s", replAndStdout)
		}
	}
}
