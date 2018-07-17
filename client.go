package main

import (
  "net"
  "fmt"
  "os"
  "os/signal"
  "syscall"
)

type Client struct {
  username string
  password string
  conn net.Conn
  uuid string 
}

func CreateClient(username string, password string, conn net.Conn) Client{
  return Client{ username, password, conn, generateUuid() }
}

func (c Client) CloseConnection() {
  fmt.Println("Closing client")
  c.conn.Close()
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
