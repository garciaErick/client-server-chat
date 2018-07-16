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

  fmt.Println("Launching server...")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", ":8085")

  SetupCloseHandlerServer(ln)
  for {

    fmt.Println("Ready to receive")

    // run loop forever (or until ctrl-c)
    conn, err := ln.Accept()

    if err != nil {
      fmt.Println("An error has ocurred, closing server")
      os.Exit(2)
    }


    go handleConnection(conn)


    fmt.Println("Closing Connection")

  }
}

func handleConnection(conn net.Conn){
  name, _ := bufio.NewReader(conn).ReadString('\n')
  name = strings.TrimRight(name, "\n")
  fmt.Println(name + " has logged in")
  for {
    // will listen for message to process ending in newline (\n)
    message, _ := bufio.NewReader(conn).ReadString('\n')


    fmt.Print("Message Received:", string(message))
    newmessage := strings.ToUpper(message)
    conn.Write([]byte(newmessage + "\n"))

    if strings.TrimRight(message, "\n") == "exit" {
      conn.Close()
      return
    }
  }

}

func SetupCloseHandlerServer(ln net.Listener) {
  c := make(chan os.Signal, 2)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-c
    fmt.Println("\r- Ctrl+C pressed in Terminal")
    fmt.Println("\r- Closing listener and exiting program")
    ln.Close()
    os.Exit(0)
  }()
}
