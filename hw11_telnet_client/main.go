package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", time.Second*10, "connection timeout")

	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage:	go-telnet --timeout=Ns host port ")
		os.Exit(1)
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	// timeout := time.Second * 2
	// address := net.JoinHostPort("127.0.0.1", "4242")

	log.SetOutput(os.Stderr)
	telnetClient := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	err := telnetClient.Connect()
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		telnetClient.Close()
		os.Stdin.Close()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		readRoutine(telnetClient, wg)
		cancel()
	}()

	wg.Add(1)
	go func() {
		writeRoutine(telnetClient, wg)
		cancel()
	}()

	wg.Wait()
	telnetClient.Close()
}

func readRoutine(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()

	err := client.Receive()
	if err != nil {
		log.Printf("Error in readRoutine: %v", err)
	}
}

func writeRoutine(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()

	err := client.Send()
	if err != nil {
		log.Printf("Error in writeRoutine: %v", err)
	}
}
