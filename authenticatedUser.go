package main

import "net"

type AuthenticatedUser struct {
  username string
  uuid     string
  conn     net.Conn
}
