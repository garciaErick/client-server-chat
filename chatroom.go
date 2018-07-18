package main

import (
  "net"
  "fmt"
  // "os"
  // "strings"
)

func main() {

  // connect to this socket
  protocol, host, port := "tcp", "localhost" , ":8085"


  fmt.Println("Starting Client connection")

  conn, _ := net.Dial(protocol, host + port)
  SetupCloseHandlerClient(conn)
  CreateClient(conn)

  fmt.Println("Connected to server ")
}

