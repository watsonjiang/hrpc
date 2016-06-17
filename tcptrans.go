package hrpc

import (
   "net"
   "fmt"
   "bytes"
   log "github.com/golang/glog"
)

//internal implementation
type TcpTrans struct {
   cfg *TransConfig
   seq  *Sequencer32
   peerReg *PeerRegistry
   timer *TimerQueue
   txQuota chan int
   listener TransListener
}

func NewTcpTrans(c *TransConfig) Trans {
   t := &TcpTrans{}
   t.cfg = c
   t.peerReg = NewPeerRegistry(t)
   t.timer = NewTimerQueue()
   t.timer.Start()
   t.seq = NewSequencer32(0)
   if c.MaxSendRate != 0 {
      go t.txQuotaLoop()
   }
   return t
}

//PeerRegistryListener
func (t *TcpTrans) OnPeerAdded(p *Peer) {
   addr := p.Addr[0]
   for i:=0;i<t.cfg.MaxConns; i++ {
      if conn, err:=net.Dial("tcp", addr);err!=nil {
         log.Errorf("Fail to connect addr[%v], err:%v", addr, err)
         continue
      }else{
         t.linkMade(p, conn)
      }
   }
}

func (t *TcpTrans) OnPeerUpdated(oldv, newv *Peer) {

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

func (t *TcpTrans) handshake(conn net.Conn) *Peer {
   peer := NewPeer(t.cfg.LocalPeerInfo)
   buf := new(bytes.Buffer)
   peer.WriteTo(buf)
   m_req := NewRequest()
   m_req.data = buf.Bytes()
   m_req.seq = 0
   m_req.mtype |= MSG_HANDSHAKE_BIT
   log.Infoln("Send handshake request", m_req)
   m_req.WriteTo(conn)

   m_rsp := NewMessage()
   if _, err:=m_rsp.ReadFrom(conn);err!=nil{
      log.Errorln("Fail to read handshake response.", err)
      conn.Close()
      return nil
   }
   if !m_rsp.IsHandshakeMsg() {
      log.Errorln("Invalid handshake response. rsp:", m_rsp)
      conn.Close()
      return nil
   }
   r_peer := &Peer{}
   if _, err:=r_peer.ReadFrom(bytes.NewBuffer(m_rsp.data));err!=nil{
      log.Errorln("Invalid handshake response. rsp:", m_rsp)
      conn.Close()
      return nil
   }
   return r_peer
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

func (t *TcpTrans) txLoop(p *Peer, cli net.Conn) {
   quotaCnt := int64(0)
   for m := range p.txChan {
      cnt, err := m.WriteTo(cli)
      if err != nil {
         //TODO: close both tx and rx loop
      }
      if t.cfg.MaxSendRate != 0{
         quotaCnt = quotaCnt - cnt
         if quotaCnt < 0 {
            for {
                quota := <-t.txQuota
                quotaCnt = quotaCnt + int64(quota * 1024)
                if quotaCnt > 0 {
                   break
                }
            }
         }
      }
   }
}

func (t *TcpTrans) rxLoop(p *Peer, cli net.Conn) {
   for {
      m := NewMessage()
      if _, err:=m.ReadFrom(cli);err!=nil{
         //TODO: close both tx and rx loop
      }
      if m.IsReqMsg() {
         m.peerId = p.Id
         t.listener.OnReqArrival(m)
      }else{
         t.listener.OnRspArrival(m)
      }
   }
}

func (t *TcpTrans) linkMade(p *Peer, cli net.Conn) {
   //TODO: handshake
   go t.rxLoop(p, cli)
   go t.rxLoop(p, cli)
}

//server
func (t *TcpTrans) Listen(addr string) error {
  go t.listenLoop()
  return nil
}

func (t *TcpTrans) linkNew(conn net.Conn) {
   r_peer := handshake(conn)
   if r_peer != nil {
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
      go t.waitForHandshake(client)
   }
}

func (t *TcpTrans) Close() {
   //TODO shutdown everything
}
