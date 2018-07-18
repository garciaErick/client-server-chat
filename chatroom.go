package main

import (
  "net"
  "fmt"
  // "os"
  // "strings"
)

func main() {

  // Using tcp protocol to ensure our messages get delivered
  protocol, host, port := "tcp", "localhost" , ":8080"


  fmt.Println("Starting Client connection")

  conn, _ := net.Dial(protocol, host + port)
  SetupCloseHandlerClient(conn)
  CreateClient(conn)

  fmt.Println("Press Ctrl+c to exit chatroom")

}

