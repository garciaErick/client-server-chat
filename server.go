package main

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	SetupCloseHandler()

	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8083")

	fmt.Println("Ready to receive")

	// run loop forever (or until ctrl-c)
	for {

		conn, _ := ln.Accept()

		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')

		if strings.TrimRight(message, "\n") == "exit" {
			break
		}

		fmt.Print("Message Received:", string(message))
		newmessage := strings.ToUpper(message)
		conn.Write([]byte(newmessage + "\n"))
	}

	fmt.Println("Closing Connection")
	ln.Close()
	os.Exit(1)
}

func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}
