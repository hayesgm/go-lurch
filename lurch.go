package main

/*
  Lurch is an error management server
*/

import (
  "github.com/hayesgm/go-lurch/server"
  "labix.org/v2/mgo"
  "log"
)

// Main runs a lease server
// First, we'll connect to a MongoDB
func main() {
  log.Println("Welcome to Lurch\n")

  session, err := mgo.Dial("localhost")
  if err != nil {
    panic(err)
  }
  defer session.Close()

  log.Println("Connected to localhost")

  server.RunLurchServer(session.DB("lurch"))
}