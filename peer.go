package hrpc

const (
   MAX_PEER_TX_QUEUE_SIZE = 1024
   MAX_PEER_RX_QUEUE_SIZE = 1024
)

type Peer struct {
   Id string
   Addr []string
   txChan  chan *Message
   rxChan  chan *Message
}

func NewPeer(info string) *Peer {
   p := &Peer{}
   if err:=json.Unmarshal(info, p);err!=nil {
      return nil
   }
   p.txChan = make(chan *Message, MAX_PEER_TX_QUEUE_SIZE)
   p.rxChan = make(chan *Message, MAX_PEER_RX_QUEUE_SIZE)
}

type PeerRegistryListener struct {
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
   if oldv, ok:=r.reg[id]; ok{
      r.reg[id] = p
      r.listener.OnPeerUpdated(oldv, p)
   }else{
      r.reg[id] = p
      r.listener.OnPeerAdded(p)
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
