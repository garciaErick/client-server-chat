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


  for {
    // listen on all interfaces
    ln, _ := net.Listen("tcp", ":8083")

    fmt.Println("Ready to receive")

    // run loop forever (or until ctrl-c)
    conn, err := ln.Accept()

    name, _ := bufio.NewReader(conn).ReadString('\n')
    name = strings.TrimRight(name, "\n")
    fmt.Println(name + " has logged in")

    for {

      if err != nil {
        fmt.Println("An error has ocurred, closing server")
        os.Exit(2)
      }

      // will listen for message to process ending in newline (\n)
      message, _ := bufio.NewReader(conn).ReadString('\n')


      fmt.Print("Message Received:", string(message))
      newmessage := strings.ToUpper(message)
      conn.Write([]byte(newmessage + "\n"))

      if strings.TrimRight(message, "\n") == "exit" {
        break
      }
    }

    fmt.Println("Closing Connection")
    ln.Close()

  }
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
