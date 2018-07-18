package main

import ( 
  "bufio"
  "strings"
  "net" 
  "log"
  "fmt"
  // "os"
  // "os/signal"
  // "syscall"
)

func main() {
  fmt.Println("Starting server...")

  ln, err := net.Listen("tcp", ":8085")

  activeUsers := make(map[string] AuthenticatedUser)
  conns  := make(chan net.Conn)
  disconnectedUsers := make(chan AuthenticatedUser)
  messages   := make(chan string)

  // SetupCloseHandlerServer(ln, activeUsers)

  fmt.Println("Server is ready to receive")

  if err != nil {
    log.Println(err.Error())
  }

  i := 0

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
      var username string

      for {
        username, _ = rd.ReadString('\n')
        username = strings.TrimRight(username, "\n")

        if UserExists(username, activeUsers){
          break
        }
      }

      uuid, _ := rd.ReadString('\n')
      uuid = strings.TrimRight(uuid, "\n") 

      message := username + " has connected\n"
      BroadcastConnection(message, activeUsers)

      activeUsers[uuid] = AuthenticatedUser {
        username: username,
        uuid: uuid,
        conn: conn,
      }

      //Read messages from connections
      go func(user AuthenticatedUser, i int){
        for {
          message, err := rd.ReadString('\n')
          message = strings.TrimRight(message, "\n")
          if err != nil {
            break
          }
          message = user.username + ": " + message
          messages <- message
        }
        disconnectedUsers <- user
      } (activeUsers[uuid], i)

    case message := <- messages:
      // Broadcast
      for _, user := range activeUsers {
        var messageToSend string
        fullMessage := strings.Split(message, ":")
        username, contents := fullMessage[0], fullMessage[1]

        if user.username == username{
          messageToSend = "(You):" + contents
        } else {
          messageToSend = username + ": " + contents

        }

        fmt.Fprintf(user.conn, messageToSend + "\n")
        user.conn.Write([]byte(messageToSend))
      }

    case disconnectedUser := <- disconnectedUsers:
      message := "User " + disconnectedUser.username + " has closed the chat"
      delete(activeUsers, disconnectedUser.uuid)

      BroadcastConnection(message, activeUsers)
    }
  }
}

func UserExists(username string, activeUsers map[string]AuthenticatedUser) (exists bool) {
  for _, user := range activeUsers{
    if user.username == username {
      return false
    }
  }
  return true
}

func BroadcastConnection (message string, activeUsers map[string]AuthenticatedUser){
  log.Printf(message)
  for _, user := range activeUsers {
    fmt.Fprintf(user.conn, message + "\n")
    user.conn.Write([]byte(message))
  }
}


// func SetupCloseHandlerServer(ln net.Listener, activeUsers map[string] AuthenticatedUser ) {
//   c := make(chan os.Signal, 2)
//   signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//   go func() {
//     <-c
//     // for conn := range activeUsers{
//     //   conn.Close()
//       // fmt.Println(k)
//     }
//     fmt.Println("\r- Ctrl+C pressed in Terminal")
//     fmt.Println("\r- Closing listener and exiting program")
//     ln.Close()
//     os.Exit(0)
//   }()
// }
