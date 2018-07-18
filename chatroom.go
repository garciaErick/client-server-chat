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


  // read in input from stdin
  // reader := bufio.NewReader(os.Stdin)

  // // Sending username
  // fmt.Print("Enter your username: ")
  // username, _ := reader.ReadString('\n')
  // fmt.Fprintf(conn, username + "\n")

  // // Sending password
  // fmt.Print("Enter your password: ")
  // password, _ := reader.ReadString('\n')
  // fmt.Fprintf(conn, password + "\n")

  // // Sending uuid
  // client := CreateClient(username, password, conn)
  // fmt.Fprintf(conn, client.uuid+ "\n")
  
  // for { 

    // fmt.Print("$_: ")
    // text, _ := reader.ReadString('\n')

    // send to socket
    // fmt.Fprintf(conn, text + "\n")

    // listen for reply
    // message, _ := bufio.NewReader(conn).ReadString('\n')
    // fmt.Print("Message from server: "+message)

    // if strings.TrimRight(text, "\n") == "exit" {
    //   fmt.Println("Closing client")
    //   break
  // }

  // if text  == "exit" {
  //   client.CloseConnection()
  // }
}

