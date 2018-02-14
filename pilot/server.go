package haedal

import (
  "net/http"
)

type Server struct {
  *router
  middlewares []Middleware
  startHandler HandlerFunc // 체인 형태 미들웨어 시작점
}

func NewServer() *Server {
  r := NewRouter()
  s := &Server{router: r}
  s.middlewares = []Middleware{
    LogHandler,
    RecoverHandler,
    StaticMiddleware("./pilot/example"), // static middleware path should be considered from cwd
  }
  return s
}

func (s *Server) Run(addr string) {
  s.startHandler = s.router.handler()

  for i := len(s.middlewares)-1; i >= 0; i-- {
    s.startHandler = s.middlewares[i](s.startHandler)
  }

  if err := http.ListenAndServe(addr, s); err != nil {
    panic(err)
  }
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  c := &Context{
    Params: make(map[string]interface{}),
    ResponseWriter: w,
    Request: r,
  }
  for k, v := range r.URL.Query() {
    c.Params[k] = v[0]
  }
  s.startHandler(c)
}
