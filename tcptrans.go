package hrpc

//internal implementation
type TcpTrans struct {
   cfg *TransConfig
   seq  *Sequencer
   peerReg PeerRegistry
   timer *TimerQueue
   txQuota chan int
   listener TransListener
}

func NewTcpTrans(c *TransConfig) *Trans {
   t := &TcpTrans{}
   t.cfg = c
   t.peerReg = NewPeerRegistry(t)
   t.timer = NewTimerQueue()
   go t.timer.Run()
   t.seq = NewSequencer(0)
   if c.MaxSendRate != 0 {
      go t.txQuotaLoop()
   }
}

//PeerRegistryListener
func (t *TcpTrans) OnPeerAdded(p *Peer) {
   addr := p.Addr[0]
   for i:=0;i<t.cfg.MaxConns; i++ {
      if conn, err:=net.Dial("tcp", addr);err!=nil {
         log.Errorf("Fail to connect addr[%v], err:%v", addr, err)
         continue
      }
      t.connectionMade(p, conn)
   }
}

func (t *TcpTrans) OnPeerUpdated(oldv, newv *Peer) {

}

//Trans
func (t *TcpTrans) Send(peerId string, m *Message) {
   t.GetTxChan(peerId) <- m   
}

func (t *TcpTrans) GetTxChan(peerId string) chan *Message {
   return t.peerReg.Get(peerId) 
}

func (t *TcpTrans) RegisteListener(l TransListener) {
   t.listener = l
}

func (t *TcpTrans) AddPeer(peerInfo string) error {
   p := NewPeer(peerInfo)
   p.txQuota = make(chan int, txQuota)
   t.peerReg.Add(peer)
   for i:=0; i<c.connNum; i++ {
      conn, err := net.Dial("tcp", c.addr)
      fmt.Println("cli make conn", c.addr, conn, err, "mode", c.mode)
      c.connMade(conn)
   }
}

func (t *TcpTrans) txQuotaLoop() {
   ch := make(chan int64, 1)
   for {
      r := t.cfg.MaxSendRate
      for i:=0;i<r;i++ {
         if len(t.txQuota) < r {
            c.txQuota<-1
         }
      }
      t.timer.Schedule(1000, ch)
      <-ch
   }
}

func (t *TcpTrans) txLoop(p *Peer, cli net.Conn) {
   var tmp_buf bytes.Buffer
   quotaCnt := 0
   for m := range p.txChan {
      tmp_buf.Reset()
      data := m.Bytes()
      binary.Write(&tmp_buf, binary.LittleEndian, int32(len(data)))
      binary.Write(&tmp_buf, binary.LittleEndian, data)
      if t.cfg.MaxSendRate != 0{
         tmpl := tmp_buf.Len()
         quotaCnt = quotaCnt - tmpl
         if quotaCnt < 0 {
            for {
                quota := <-t.txQuota
                quotaCnt = quotaCnt + quota * 1024
                if quotaCnt > 0 {
                   break
                }
            }
         }
      }
      cli.Write(tmp_buf.Bytes())
   }
}

func (t *TcpTrans) rxLoop(p *Peer, cli net.Conn) {
   for {
      var l int32
      binary.Read(cli, binary.LittleEndian, &l)
      data := make([]byte, l)
      binary.Read(cli, binary.LittleEndian, &data)
      m := decodeMessage(data)
      if m.isReqMsg() {
         t.listener.OnReqArrival(m)
      }else{
         t.listener.OnRspArrival(m)
      }
   }
}

func (t *TcpTrans) connectionMade(p *Peer, cli net.Conn) {
   //TODO: handshake
   go t.rxLoop(p, cli)
   go t.rxLoop(p, cli)
}

//server
func (t *TcpTrans) Listen(addr string) error {
  go t.listenLoop()
}

func (t *TcpTrans) listenLoop() {
   listener, err := net.Listen("tcp", s.addr)
   fmt.Println("svr listen on ", s.addr, "err:", err)
   for {
      client, _ := listener.Accept()
      fmt.Println("svr got conn", client)
      go t.connMade(client)
   }
}

func (t *TcpTrans) Close() {
   //TODO shutdown everything
}
