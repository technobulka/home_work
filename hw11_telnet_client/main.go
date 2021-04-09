package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Fatal("usage: go-telnet [--timeout=10s] host port")
	}
	log.SetOutput(os.Stderr)

	shutdown := make(chan bool, 2)
	pr, pw := io.Pipe()
	go handleInput(pw, shutdown)

	addr := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	client := NewTelnetClient(addr, timeout, pr, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to", addr)
	go receiveRoutine(client, shutdown)
	go sendRoutine(client, shutdown)
	handleInterrupt(client, shutdown)
}

func handleInput(w io.Writer, shutdown chan bool) {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		_, err := w.Write([]byte(input.Text() + "\n"))
		if err != nil {
			log.Println(err)
			shutdown <- true
			return
		}
	}

	if err := input.Err(); err != nil {
		log.Println(err)
		shutdown <- true
		return
	}

	shutdown <- true
}

func receiveRoutine(client TelnetClient, shutdown chan bool) {
	if err := client.Receive(); err != nil {
		log.Println("error in receiving:", err)
		shutdown <- true
	}
}

func sendRoutine(client TelnetClient, shutdown chan bool) {
	if err := client.Send(); err != nil {
		log.Println("error in sending:", err)
		shutdown <- true
	}
}

func handleInterrupt(client TelnetClient, shutdown chan bool) {
	signalCh := make(chan os.Signal, 1)
	defer close(signalCh)

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalCh)

	select {
	case s := <-signalCh:
		log.Println("received:", s)
	case <-shutdown:
	}

	if err := client.Close(); err != nil {
		log.Println(err)
	}

	log.Println("Connection was closed")
}
