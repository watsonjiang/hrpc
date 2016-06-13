package hrpc

import (
   "sync"
   "errors"
   "encoding/json"
   "io"
   "encoding/binary"
)

const (
   MAX_PEER_TX_QUEUE_SIZE = 1024
   MAX_PEER_RX_QUEUE_SIZE = 1024
)

var ERR_INVALID_PEER = errors.New("Invalid Peer info")

type Peer struct {
   Id string
   Addr []string
   txChan  chan *Message
}

func NewPeer(info string) *Peer {
   p := &Peer{}
   if err:=json.Unmarshal([]byte(info), p);err!=nil {
      return nil
   }
   p.txChan = make(chan *Message, MAX_PEER_TX_QUEUE_SIZE)
   return p
}

func (p *Peer) WriteTo(w io.Writer) (int64, error) {
   cnt := int64(0)
   if err:=binary.Write(w, binary.LittleEndian, int8(len(p.Id)));err!=nil{
      return 0, err
   }
   cnt += int64(1)
   if n, err:=w.Write([]byte(p.Id));err!=nil {
      cnt += int64(n)
      return cnt, err
   }
   cnt += int64(len(p.Id))
   if err:=binary.Write(w, binary.LittleEndian, int8(len(p.Addr)));err!=nil{
      return cnt, err
   }
   cnt += int64(1)
   for _, addr := range p.Addr {
      if err:=binary.Write(w, binary.LittleEndian, int8(len(addr)));err!=nil {
         return cnt, err
      }
      cnt += int64(1)
      if n, err:=w.Write([]byte(addr));err!=nil{
         cnt += int64(n)
         return cnt, err
      }
      cnt += int64(len(addr))
   }
   return cnt, nil
}

func (p *Peer) ReadFrom(r io.Reader) (int64, error) {
   cnt := int64(0)
   var lid int8
   if err:=binary.Read(r, binary.LittleEndian, &lid);err!=nil{
      return 0, err
   }
   cnt += int64(1)
   var id_buf = make([]byte, lid)
   if n, err:=r.Read(id_buf);err!=nil{
      cnt += int64(n)
      return cnt, err
   }
   cnt += int64(lid)
   p.Id = string(id_buf)
   var n_addr int8
   if err:=binary.Read(r, binary.LittleEndian, &n_addr);err!=nil{
      return cnt, err
   }
   cnt += 1
   for i:=0; i<int(n_addr); i++ {
      var l_addr int8
      if err:=binary.Read(r, binary.LittleEndian, &l_addr);err!=nil{
         return cnt, err
      }
      cnt += 1
      buf := make([]byte, l_addr)
      if n, err:=r.Read(buf);err!=nil {
         cnt += int64(n)
         return cnt, err
      }
      cnt += int64(l_addr)
      p.Addr = append(p.Addr, string(buf))
   }
   return cnt, nil
}

type PeerRegistryListener interface {
   OnPeerAdded(p *Peer)
   OnPeerUpdated(oldv, newv *Peer)
}

type PeerRegistry struct {
   reg map [string] *Peer
   lock sync.Mutex
   listener PeerRegistryListener
}

func NewPeerRegistry(l PeerRegistryListener) *PeerRegistry {
   return &PeerRegistry{
      reg:make(map[string]*Peer),
      listener:l,
   }
}

func (r *PeerRegistry) Put(p *Peer) {
   r.lock.Lock()
   defer r.lock.Unlock()
   if oldv, ok:=r.reg[p.Id]; ok{
      r.reg[p.Id] = p
      if r.listener != nil {
         r.listener.OnPeerUpdated(oldv, p)
      }
   }else{
      r.reg[p.Id] = p
      if r.listener != nil {
         r.listener.OnPeerAdded(p)
      }
   }
}

func (r *PeerRegistry) Get(id string) *Peer {
   r.lock.Lock()
   defer r.lock.Unlock()
   if p, ok := r.reg[id];ok {
      return p
   }
   return nil
}

