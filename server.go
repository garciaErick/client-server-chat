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

var clients_MAP = make(map[int]clients_table)

func main() {

  fmt.Println("Launching server...")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", ":8085")

  SetupCloseHandlerServer(ln)

  fmt.Println("Ready to receive connections")
  fmt.Println()

  for {

    // run loop forever (or until ctrl-c)
    conn, err := ln.Accept()

    if err != nil {
      fmt.Println("An error has ocurred, closing server")
      os.Exit(2)
    }

    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn){
  username, _ := bufio.NewReader(conn).ReadString('\n')
  username = strings.TrimRight(username, "\n")
  fmt.Println("Received user info " + username )

  password, _ := bufio.NewReader(conn).ReadString('\n')
  password = strings.TrimRight(password, "\n")
  fmt.Println("Received user info " + password)

  uuid, _ := bufio.NewReader(conn).ReadString('\n')
  uuid = strings.TrimRight(uuid, "\n")
  fmt.Println("Received user info " + uuid)

  client := Client{username, password, conn, uuid}
  clients_MAP[uuid] = client

  fmt.Println(client.username + " has logged in")
  for {
    // will listen for message to process ending in newline (\n)
    message, _ := bufio.NewReader(conn).ReadString('\n')


    fmt.Print("Message Received:", string(message))
    newmessage := strings.ToUpper(message)
    conn.Write([]byte(newmessage + "\n"))

    if strings.TrimRight(message, "\n") == "exit" {
      fmt.Println(client.username + " has disconected")
      conn.Close()
      fmt.Println("Closing Connection")
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
