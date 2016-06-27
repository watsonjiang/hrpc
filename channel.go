package hrpc

type Channel struct {

}


type ChannelRegistry struct {
   reg map [string] *Peer
   lock sync.Mutex
}

func NewChannelRegistry() *ChannelRegistry {
   return &ChannelRegistry{
      reg:make(map[string]*Peer),
   }
}

func (r *ChannelRegistry) Put(p *Peer) {
   r.lock.Lock()
   defer r.lock.Unlock()
   r.reg[p.Id] = p
}

func (r *ChannelRegistry) Get(id string) *Peer {
   r.lock.Lock()
   defer r.lock.Unlock()
   if p, ok := r.reg[id];ok {
      return p
   }
   return nil
}

type func JobT(p *Peer) error
// map all elem peers using function job
func (r *ChannelRegistry) filter(job JobT) error {
   r.lock.Lock()
   defer r.lock.Unlock()
   for _, p := range r.reg {
      job(p)
   }
   return nil
}

func (r *ChannelRegistry) doWith(id, job jobT) error {
   return nil
}
