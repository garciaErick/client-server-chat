package main

import (
  "net"
  "fmt"
  "bufio"
  "os"
  "strings"
)

func main() {

  // connect to this socket
  conn, _ := net.Dial("tcp", "127.0.0.1:8085")
  SetupCloseHandlerClient(conn)

  fmt.Println("Starting Client connection")

  // read in input from stdin
  reader := bufio.NewReader(os.Stdin)

  fmt.Print("Enter your username: ")
  username, _ := reader.ReadString('\n')
  fmt.Fprintf(conn, username + "\n")
  fmt.Print("Enter your password: ")
  password, _ := reader.ReadString('\n')
  fmt.Fprintf(conn, password + "\n")

  client := CreateClient(username, password, conn)
  fmt.Println("Sent user info " + username + " " + password + " " + client.uuid)
  fmt.Fprintf(conn, client.uuid+ "\n")

  for { 

    fmt.Print("$_ ")
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
      client.CloseConnection()
    }
  }
}
