package hrpc

type Server struct {
   addr string
}

func (s *Server) Start(){
   go s.listenLoop()
}

func (s *Server) handleRequest(req *Request) *Response {
   rsp := &Response{}
   rsp.seq = req.seq
   rsp.data = []byte("world")
   return rsp
}

func (s *Server) listenLoop() {
   listener, err := net.Listen("tcp", s.addr)
   fmt.Println("svr listen on ", s.addr, "err:", err)
   for {
      client, _ := listener.Accept()
      fmt.Println("svr got conn", client)
      go s.connMade(client)
   }
}


