package hrpc

import (
   "net"
   "bytes"
   log "github.com/golang/glog"
   "time"
)

type Link struct {
   status int
   listener TransListener
   flowCtlEnabled bool
   txQuota chan int
   txChan  chan *Message
   stop   int
}

func sendHandshake(p *Peer, c net.Conn) error {
   buf := new(bytes.Buffer)
   p.WriteTo(buf)
   m_req := NewRequest()
   m_req.data = buf.Bytes()
   m_req.seq = 0
   m_req.mtype |= MSG_HANDSHAKE_BIT
   log.Infoln("Send handshake message", m_req)
   c.SetWriteDeadline(time.Now().Add(HANDSHAKE_TIMEOUT))
   if _, err := m_req.WriteTo(c);err!=nil{
      log.Errorln("Fail to send handshake message.", err)
      return err
   }
   c.SetWriteDeadline(time.Time{})
   return nil
}

var HANDSHAKE_TIMEOUT = 1 * time.Second

func readHandshake(c net.Conn) (*Peer, error){
   m := NewMessage()
   c.SetReadDeadline(time.Now().Add(HANDSHAKE_TIMEOUT))
   if _, err:=m.ReadFrom(c);err!=nil{
      log.Errorln("Fail to read handshake message.", err)
      return nil, err
   }
   if !m.IsHandshakeMsg() {
      log.Errorln("Invalid handshake message.", m)
      return nil, ERR_INVALID_HANDSHAKE_MSG
   }
   r_peer := &Peer{}
   if _, err:=r_peer.ReadFrom(bytes.NewBuffer(m.data));err!=nil{
      log.Errorln("Invalid handshake message.", m, err)
      return nil, err
   }
   c.SetReadDeadline(time.Time{})
   return r_peer, nil
}


func handshake(l_peer *Peer, conn net.Conn) (*Peer, error) {
   if err:=sendHandshake(l_peer, conn);err!=nil{
      return nil, err
   }
   m := NewMessage()
   if _, err:=m.ReadFrom(conn);err!=nil{
      log.Errorln("Fail to read handshake message.", err)
      return nil, err
   }
   if !m.IsHandshakeMsg() {
      log.Errorln("Invalid handshake message.", m)
      return nil, ERR_INVALID_HANDSHAKE_MSG
   }
   r_peer := &Peer{}
   if _, err:=r_peer.ReadFrom(bytes.NewBuffer(m.data));err!=nil{
      log.Errorln("Invalid handshake message.", m, err)
      return nil, err
   }
   return r_peer, nil
}

func (l *Link) txLoop(p *Peer, cli net.Conn) {
   quotaCnt := int64(0)
   for m := range l.txChan {
      cnt, err := m.WriteTo(cli)
      if err != nil {
         //TODO: close both tx and rx loop
      }
      if l.flowCtlEnabled{
         quotaCnt = quotaCnt - cnt
         if quotaCnt < 0 {
            for {
                quota := <-l.txQuota
                quotaCnt = quotaCnt + int64(quota * 1024)
                if quotaCnt > 0 {
                   break
                }
            }
         }
      }
   }
}

func (l *Link) rxLoop(p *Peer, cli net.Conn) {
   for {
      m := NewMessage()
      if _, err:=m.ReadFrom(cli);err!=nil{
         //TODO: close both tx and rx loop
      }
      if m.IsReqMsg() {
         m.peerId = p.Id
         l.listener.OnReqArrival(m)
      }else{
         l.listener.OnRspArrival(m)
      }
   }
}


