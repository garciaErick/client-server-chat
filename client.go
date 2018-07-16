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

type client struct {
  name string
  id int
}


func main() {

  // connect to this socket
  conn, _ := net.Dial("tcp", "127.0.0.1:8085")
  SetupCloseHandlerClient(conn)

  generateUuid()

  fmt.Println("Starting Client")

  // read in input from stdin
  reader := bufio.NewReader(os.Stdin)
  fmt.Println("Enter your name")

  name, _ := reader.ReadString('\n')
  fmt.Fprintf(conn, name + "\n")

  for { 

    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')


    // send to socket
    fmt.Fprintf(conn, text + "\n")

    // listen for reply
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Print("Message from server: "+message)

    if strings.TrimRight(text, "\n") == "exit" {
      fmt.Println("Closing client")
      break
    }
    if text  == "exit" {
      fmt.Println("Closing client")
      break
    }
  }

  conn.Close()
  os.Exit(1)
}

func SetupCloseHandlerClient(conn net.Conn) {
  c := make(chan os.Signal, 2)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-c

    //Sending close connection string to server
    fmt.Fprintf(conn, "exit" + "\n")
    fmt.Println("\r- Ctrl+C pressed in Terminal, closing connection")
    os.Exit(0)
  }()
}
