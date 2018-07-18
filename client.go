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
  uuid string 
  username string
  incoming chan string
  outgoing chan string
  reader *bufio.Reader
  writer *bufio.Writer
  conn net.Conn
}

func CreateClient(conn net.Conn) Client{
  writer := bufio.NewWriter(conn)
  reader := bufio.NewReader(conn)

  username, uuid := CreateCredentials()

  client := Client{
    uuid: uuid,
    username: username,
    reader:   reader,
    writer:   writer,
    incoming: make(chan string),
    outgoing: make(chan string),
    conn:     conn,
  }

  client.Listen()

  return client
}

func CreateCredentials() (string, string){
  reader := bufio.NewReader(os.Stdin)

  fmt.Print("Enter your username:")
  username, _ := reader.ReadString('\n')
  username = strings.TrimRight(username, "\n")

  uuid := GenerateUuid()

  return username, uuid
}


func (client *Client) Read() {
  for {
    m, err := bufio.NewReader(client.conn).ReadString('\n')
    m = strings.TrimRight(m, "\n")
    if err != nil {
      break
    }
    fmt.Println(m)
    fmt.Print(">_ ")
  }
}

func (client *Client) Write() {
  for {
    fmt.Print(">_ ")
    message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    message = strings.TrimRight(message, "\n")
    fmt.Fprintf(client.conn, message + "\n")
  }
}

func (client *Client) Listen() {
  fmt.Fprintf(client.conn, client.username + "\n")
  fmt.Fprintf(client.conn, client.uuid + "\n")
  

  go client.Read()
  client.Write()
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
