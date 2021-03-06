package hrpc

import (
   "net"
   "fmt"
   "time"
   log "github.com/golang/glog"
)

//internal implementation
type TcpTrans struct {
   cfg *TransConfig
   seq  *Sequencer32
   chReg *PeerRegistry
   timer *TimerQueue
   txQuota chan int
   listener TransListener
}

func NewTcpTrans(c *TransConfig) Trans {
   t := &TcpTrans{}
   t.cfg = c
   t.chReg = NewChannelRegistry()
   t.timer = NewTimerQueue()
   t.timer.Start()
   t.seq = NewSequencer32(0)
   if c.MaxSendRate != 0 {
      go t.txQuotaLoop()
   }
   return t
}

//background routine, create/remove link as needed.
func (t *TcpTrans) linkMonitorLoop() {
   job := func(p *Peer) error{
      if len(p.links) == 0 {
         //create link 
	 l_peer := NewPeer(t.cfg.LocalPeerInfo)
	 link := NewLink(l_peer, p)
         if link != nil{
	    link.listener = t.listener
	    link.flowCtlEnabled = false
	    link.txChan = p.txChan
	    link.Start()
            p.links = append(p.links, link)
	 }
      }
   }
   for {
      time.Sleep(3*time.Second)
      t.peerReg.Map(job)
   }
}

//Trans
func (t *TcpTrans) Send(m *Message) {
   t.getTxChan(m.peerId) <- m
}

func (t *TcpTrans) getTxChan(peerId string) chan *Message {
   return t.peerReg.Get(peerId).txChan
}

func (t *TcpTrans) RegisteListener(l TransListener) {
   t.listener = l
}

func (t *TcpTrans) AddPeer(peerInfo string) error {
   p := NewPeer(peerInfo)
   t.peerReg.Put(p)
   for i:=0; i<t.cfg.MaxConns; i++ {
      if log.V(1) {
         log.V(1).Infoln("cli make conn", p.Addr)
      }
      if conn, err:=net.Dial("tcp", p.Addr[0]);err!=nil{
         log.Errorln("Fail to dial", p.Addr[0], err)
      }else{
         go t.linkNew(conn)
      }
   }
   return nil
}

func (t *TcpTrans) txQuotaLoop() {
   ch := make(chan int64, 1)
   for {
      r := t.cfg.MaxSendRate
      for i:=0;i<r;i++ {
         if len(t.txQuota) < r {
            t.txQuota<-1
         }
      }
      t.timer.Schedule(1000, ch)
      <-ch
   }
}

func (t *TcpTrans) linkMade(p *Peer, cli net.Conn) {
   //TODO: handshake
   //go t.rxLoop(p, cli)
   //go t.rxLoop(p, cli)
}

//server
func (t *TcpTrans) Listen(addr string) error {
  go t.listenLoop()
  return nil
}

func (t *TcpTrans) linkNew(conn net.Conn) {
   l_peer := NewPeer(t.cfg.LocalPeerInfo)
   if r_peer, err:=handshake(l_peer, conn);err!=nil{
      log.Errorln("Handhake failed. close connection.")
      conn.Close()
      return
   }else{
      t.linkMade(r_peer, conn)
   }
}

func (t *TcpTrans) listenLoop() {
   p := NewPeer(t.cfg.LocalPeerInfo)
   addr := p.Addr[0]
   listener, err := net.Listen("tcp", addr)
   fmt.Println("svr listen on ", addr, "err:", err)
   for {
      client, _ := listener.Accept()
      fmt.Println("svr got conn", client)
      go t.linkNew(client)
   }
}

func (t *TcpTrans) Close() {
   //TODO shutdown everything
}
