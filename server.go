package main

import ( 
  "bufio"
  "strings"
  "net" 
  "log"
  "fmt"
  "os"
  "os/signal"
  "syscall"
)

func main() {
  fmt.Println("Starting server...")

  // Using tcp protocol to ensure our messages get delivered
  protocol, host, port := "tcp", "localhost" , ":8080"
  ln, err := net.Listen(protocol, host + port) 

  activeUsers       := make(map[string] AuthenticatedUser)
  conns             := make(chan net.Conn)
  disconnectedUsers := make(chan AuthenticatedUser)
  messages          := make(chan string)

  SetupCloseHandlerServer(ln, activeUsers)
  fmt.Println("Server started at "  + host + port)
  fmt.Println("Press Ctrl+c to shutdown server")

  if err != nil {
    log.Println(err.Error())
  }

  // Handle incomming connections
  go func() {
    for {
      conn, err := ln.Accept()
      if err != nil {
        log.Println(err.Error())
      }
      conns <- conn
    }
  }()

  for {
    select {
      // Read incoming connections
    case conn := <-conns:
      rd := bufio.NewReader(conn)

      // Validating incoming user
      uuid, _ := rd.ReadString('\n')
      uuid     = strings.TrimRight(uuid, "\n")

      var username string

      // Checking if username is unique 
      for {
        username, _ = rd.ReadString('\n')
        username    = strings.TrimRight(username, "\n")

        message := ""

        // Prompt the client to use a different username
        if UserExists(username, activeUsers){
          log.Printf("Connection '%v' tried to log in with a duplicate username", uuid)
          message = ""
          message = "Username taken, try again: "
          fmt.Fprintf(conn, message + "\n")
          conn.Write([]byte(message))
          continue
        } else if ContainsIllegalCharacters(username){
          log.Printf("Connection '%v' tried to log in with username containing illegal chars", uuid)
          message = ""
          message = "Username contains illegal characters, try again: "
          fmt.Fprintf(conn, message + "\n")
          conn.Write([]byte(message))
          continue
        } else {
          // Successfully authenticated user
          message = "Connection stablished with the chat server"
          fmt.Fprintf(conn, message + "\n")
          conn.Write([]byte(message))
          break
        }
      }

      // Broadcasting new user to active connections
      message := username + " has connected\n"
      BroadcastConnection(message, activeUsers)

      // Add new user to the Authenticated Users
      activeUsers[uuid] = AuthenticatedUser {
        username: username,
        uuid:     uuid,
        conn:     conn,
      }

      // Read messages from connections authenticated users
      go func(user AuthenticatedUser){
        for {
          // Declaring empty message to clean any leftover buffers
          message       = ""
          message, err := rd.ReadString('\n')
          message       = strings.TrimRight(message, "\n")

          if err != nil {
            break
          }
          message = user.username + ": " + message

          // Handle incoming messages
          messages <- message
        }

        // Handle when user disconnects from the chat
        disconnectedUsers <- user
      } (activeUsers[uuid])


      // Broadcast incoming message to Authenticated Users
    case message := <- messages:
      for _, user := range activeUsers {
        messageToSend      := ""
        fullMessage        := strings.Split(message, ":")
        username, contents := fullMessage[0], fullMessage[1]

        // Checking if sending message to the original sender
        if user.username == username{
          messageToSend = "(You):" + contents
        } else {
          messageToSend = username + ": " + contents
        }

        fmt.Fprintf(user.conn, messageToSend + "\n")
        user.conn.Write([]byte(messageToSend))
      }

      // Broadcast disconnecting user and removing from Authenticated Users
    case disconnectedUser := <- disconnectedUsers:
      message := ""
      if disconnectedUser.username == "" {
        message = "Unknown user has closed the chat"
      } else {
        message = "User " + disconnectedUser.username + " has closed the chat"
      }
      delete(activeUsers, disconnectedUser.uuid)

      BroadcastConnection(message, activeUsers)
    }
  }
}

func UserExists(username string, activeUsers map[string]AuthenticatedUser) (bool) {
  for _, user := range activeUsers{
    if user.username == username {
      return true
    }
  }
  return false
}

func ContainsIllegalCharacters(username string) (bool) {
  return strings.Contains(username, ":")
}

func BroadcastConnection (message string, activeUsers map[string]AuthenticatedUser){
  log.Printf(message)
  for _, user := range activeUsers {
    fmt.Fprintf(user.conn, "From server: " + message + "\n")
    user.conn.Write([]byte(message))
  }
}


func SetupCloseHandlerServer(ln net.Listener, activeUsers map[string] AuthenticatedUser ) {
  c := make(chan os.Signal, 2)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-c
    fmt.Println()
    log.Printf("Ctrl+C pressed in Terminal")
    log.Printf("Closing listener, connections, and exiting program")

    BroadcastConnection("Server is shutting down, closing connections", activeUsers)

    for _, user := range activeUsers{
      user.conn.Close()
      log.Printf("Closing connetion with user " + user.username)
    }

    log.Printf("Closing listener, Goodbye!")
    ln.Close()
    os.Exit(0)
  }()
}
