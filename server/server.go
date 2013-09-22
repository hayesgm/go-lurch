package server

import (
  "net/http"
  "encoding/json"
  "log"
  "labix.org/v2/mgo"
  "strings"
  "errors"
  "github.com/hayesgm/go-lurch/lurch"
)

type ErrorResponse struct {
  Error string
}

type TrackExceptionResponse struct {
  Track *lurch.Exception
}

// Track Error will track the instance of an error
//
// Signature:
// <-
// name:: Name of the error (e.g. NotFoundError)
// message:: Message associated with error (e.g. "Unable to find /users/5")
// backtrace:: Optional, comma-separated list which is backtrace of error
// context:: JSON-encoded context object for storing arbitrary data
// -> (json)
// { Track: '<tracked error information>'}
func HandleTrackException(db *mgo.Database, w http.ResponseWriter, req *http.Request) {
  var context interface{}
  var backtrace []string

  // First, we'll grab the parameters
  name := req.FormValue("name")
  message := req.FormValue("message")
  
  if len(name) == 0 || len(message) == 0 {
    panic(errors.New("Missing required name or message"))
  }

  backtraceV := req.FormValue("backtrace")
  if len(backtraceV) > 0 {
    backtrace = strings.Split(backtraceV,"\n")
  }

  contextV := req.FormValue("context")

  if len(contextV) > 0 {
    err := json.Unmarshal([]byte(contextV), &context)
    if err != nil {
      panic(err)
    }
  }
  
  // Now we'll create a shell object for the error
  exception, err := lurch.NewException(name, message, backtrace)
  if err != nil {
    panic(err)
  }

  // And a shell trace-context
  trace, err := lurch.NewTrace(message, backtrace, context)
  if err != nil {
    panic(err)
  }

  // Finally, we'll upsert this into the database
  err = exception.Upsert(trace, db)
  if err != nil {
    panic(err)
  }

  exception, err = lurch.FindException(exception.UniqueId, db) // Reload changes
  if err != nil {
    panic(err)
  }

  resp := TrackExceptionResponse{Track: exception}
  respBytes, err := json.Marshal(&resp)
  if err != nil {
    panic(err)
  }
  
  if _, err := w.Write(respBytes); err != nil {
    panic(err)
  }
}

func RunLurchServer(db *mgo.Database) {

  handleWithDatabase := func(f func(db *mgo.Database, w http.ResponseWriter, req *http.Request)) (func(w http.ResponseWriter, req *http.Request)) {
    return func(w http.ResponseWriter, req *http.Request) {
      defer func() {
        return

        if e := recover(); e != nil {
          resp := ErrorResponse{Error: e.(error).Error()}
          respBytes, err := json.Marshal(&resp)
          if err != nil {
            panic(err)
          }
          
          if _, err := w.Write(respBytes); err != nil {
            panic(err)
          }
        }
      }()

      f(db, w, req)
    }
  }

  http.Handle("/track-exception", http.HandlerFunc(handleWithDatabase(HandleTrackException)))
  
  log.Println("Running Lurch Server...")

  err := http.ListenAndServe("localhost:9119", nil)
  if err != nil {
    log.Fatal("ListenAndServer:",err)
  }
}