package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()
	host := flag.Arg(0)
	port := flag.Arg(1)

	if len(flag.Args()) < 2 {
		log.Fatal("usage: go-telnet [--timeout=10s] host port")
	}

	log.SetOutput(os.Stderr)

	addr := net.JoinHostPort(host, port)
	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	handleConnection(ctx, client)

	go readRoutine(ctx, client, cancel)
	go writeRoutine(ctx, client, cancel)

	err = client.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("close telnet client")
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}

	log.Println("Bye, client!")
}

func handleConnection(ctx context.Context, client TelnetClient) {
	defer client.Close()
	//client.Write([]byte(fmt.Sprintf("Welcome to %s, friend from %s\n", conn.LocalAddr(), conn.RemoteAddr())))

	//scanner := bufio.NewScanner(conn)
	//for scanner.Scan() {
	//	text := scanner.Text()
	//	log.Printf("RECEIVED: %s", text)
	//	if text == "quit" || text == "exit" {
	//		break
	//	}
	//
	//	conn.Write([]byte(fmt.Sprintf("I have received '%s'\n", text)))
	//}
	//
	//if err := scanner.Err(); err != nil {
	//	log.Printf("Error happend on connection with %s: %v", conn.RemoteAddr(), err)
	//}
	//
	//log.Printf("Closing connection with %s", conn.RemoteAddr())

}

func readRoutine(ctx context.Context, client TelnetClient, cancelFunc context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			cancelFunc()
			return
		default:
			if err := client.Receive(); err != nil {
				log.Println("Receive error:", err)
				cancelFunc()
				return
			}
		}
	}
}

func writeRoutine(ctx context.Context, client TelnetClient, cancelFunc context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			cancelFunc()
			return
		default:
			if err := client.Send(); err != nil {
				log.Println("Send error:", err)
				cancelFunc()
				return
			}
		}
	}
}
