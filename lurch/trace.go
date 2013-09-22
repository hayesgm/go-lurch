package lurch

type Trace struct {
  Message string
  Backtrace []string
  Context interface{}
}

func NewTrace(message string, backtrace []string, context interface{}) (trace *Trace, err error) {
  trace = &Trace{Message: message, Backtrace: backtrace, Context: context}
  
  return
}
