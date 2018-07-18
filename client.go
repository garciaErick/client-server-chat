package main

import (
  "net"
  "fmt"
  "bufio"
  "os"
  "os/signal"
  "syscall"
  "strings"
)

type Client struct {
  uuid     string
  username string
  incoming chan string
  outgoing chan string
  reader   *bufio.Reader
  writer   *bufio.Writer
  conn     net.Conn
}

func CreateClient(conn net.Conn) Client{
  // Reader and writer for communication
  writer := bufio.NewWriter(conn)
  reader := bufio.NewReader(conn)

  username, uuid := CreateCredentials()

  client := Client{
    uuid:     uuid,
    username: username,
    reader:   reader,
    writer:   writer,
    incoming: make(chan string),
    outgoing: make(chan string),
    conn:     conn,
  }

  // Start communication
  client.Listen()

  return client
}

func CreateCredentials() (string, string){
  reader := bufio.NewReader(os.Stdin)

  fmt.Print("Enter your username: ")
  username, _ := reader.ReadString('\n')
  username     = strings.TrimRight(username, "\n")

  // Unique token generated each time
  uuid := GenerateUuid()

  return username, uuid
}


func (client *Client) Read() {
  for {
    m, err := bufio.NewReader(client.conn).ReadString('\n')
    m       = strings.TrimRight(m, "\n")
    
    // If error or if server is shutdown close connection
    if err != nil  || m == "From server: Server is shutting down, closing connections"{
      fmt.Println(m)
      fmt.Println("Goodbye!")
      client.conn.Close()
      os.Exit(1)
    }
    fmt.Println(m)
    fmt.Print(">_ ")
  }
}

func (client *Client) Write() {
  for {
    fmt.Print(">_ ")

    message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    message     = strings.TrimRight(message, "\n")

    fmt.Fprintf(client.conn, message + "\n")
  }
}

func (client *Client) Listen() {
  client.LogIn()

  go client.Read()
  client.Write()
}

func (client *Client) LogIn() {
  fmt.Fprintf(client.conn, client.uuid + "\n") 

  // Repeat until successfully authenticated with server
  for {
    fmt.Fprintf(client.conn, client.username + "\n")

    serverResponse, _ := bufio.NewReader(client.conn).ReadString('\n')
    serverResponse     = strings.TrimRight(serverResponse, "\n")

    if serverResponse == "Connection stablished with the chat server" {
      fmt.Println(serverResponse)
      fmt.Println("Press Ctrl+c to exit chatroom")

      break
    } else { // Username is not unique try again
      fmt.Print(serverResponse)
      newUsername, _  := bufio.NewReader(os.Stdin).ReadString('\n')
      newUsername      = strings.TrimRight(newUsername, "\n")
      client.username  = newUsername

      fmt.Fprintf(client.conn, newUsername + "\n")
    }
  } 
}

func SetupCloseHandlerClient(conn net.Conn) {
  c := make(chan os.Signal, 2)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-c
    //Sending close connection string to server
    fmt.Println("\r- Ctrl+C pressed in Terminal, closing connection")
    conn.Close()

    os.Exit(0)
  }()
}
