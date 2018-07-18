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

  ln, err := net.Listen("tcp", ":8085")

  aconns := make(map[net.Conn]int)
  conns  := make(chan net.Conn)
  dconns := make(chan net.Conn)
  msgs   := make(chan string)

  SetupCloseHandlerServer(ln, aconns)

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
      aconns[conn] = i
      i++
      rd := bufio.NewReader(conn)
      username, _ := rd.ReadString('\n')
      username = strings.TrimRight(username, "\n")

      password, _ := rd.ReadString('\n')
      password = strings.TrimRight(password, "\n")

      fmt.Printf("%v has connected\n", username)

      //Read messages from connections
      go func(conn net.Conn, i int){
        for {
          m, err := rd.ReadString('\n')
          m = strings.TrimRight(m, "\n")
          if err != nil {
            break
          }
          msgs <- fmt.Sprintf("Client %v: %v", i, m)
        }
        dconns <- conn
      } (conn, i)

    case msg := <- msgs:
      // Broadcast
      for conn := range aconns {
        // fmt.Println(msg)
        fmt.Fprintf(conn, msg + "\n")
        conn.Write([]byte(msg))
      }

    case dconn := <- dconns:
      log.Printf("Client %v is gone\n", aconns[dconn] + 1)
      delete(aconns, dconn)
    }
  }
}


func SetupCloseHandlerServer(ln net.Listener, aconns map[net.Conn] int) {
  c := make(chan os.Signal, 2)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-c
    for conn := range aconns{
      conn.Close()
      // fmt.Println(k)
    }
    fmt.Println("\r- Ctrl+C pressed in Terminal")
    fmt.Println("\r- Closing listener and exiting program")
    ln.Close()
    os.Exit(0)
  }()
}
