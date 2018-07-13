package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"strings"
	"os/signal"
	"syscall"
)

func main() {

	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8083")

	fmt.Println("Starting Client")
	fmt.Println()

	for { 
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')

		if strings.TrimRight(text, "\n") == "exit" {
			break
		}

		// send to socket
		fmt.Fprintf(conn, text + "\n")

		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)
	}
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
