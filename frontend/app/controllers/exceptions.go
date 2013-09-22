package controllers

import (
  "github.com/robfig/revel"
  "github.com/hayesgm/go-lurch/lurch"
  "labix.org/v2/mgo"
  "log"
)

type Exceptions struct {
  *revel.Controller
}

func (c Exceptions) Index() revel.Result {
  session, err := mgo.Dial("localhost")
  if err != nil {
    panic(err)
  }
  defer session.Close()

  exceptions, err := lurch.GetExceptions(session.DB("lurch"))
  if err != nil {
    panic(err)
  }
  log.Printf("Exceptions: %#v - %d\n", exceptions[0], len(exceptions))
  return c.Render(exceptions)
}
