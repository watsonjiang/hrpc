package hrpc

type RpcHandler func(req []byte) []byte

var ERR_RPC_TIMEOUT = errors.New("RPC call timeout")

type RpcRegistryEntry struct{
   req *Message
   rsp *Message
   done chan int
}

type RpcRegistry struct {
   reg map[int32] *RpcRegistryEntry
   lock sync.Mutex
}

func NewRpcRegistry() *SyncReqRegistry {
   return &Registry{
      reg:make(map[int32] *RegistryEntry),
   }
}

func (r *RpcRegistry) add(seq int32, e *RpcRegistryEntry) {
   r.lock.Lock()
   defer r.lock.Unlock()
   r.reg[seq] = e
}

func (r *RpcRegistry) del(seq int32) *RpcRegistryEntry {
   r.lock.Lock()
   defer r.lock.Unlock()
   if e, ok:= r.reg[seq]; ok {
      delete(r, seq)
      return e
   }
   return nil
}

type RpcAdapter struct {
   trans Trans
   seq *Sequencer
   rpcReg *RpcRegistry
   
}

func (a *RpcAdapter) RegisterCall() chan int {
   
}

func (a *RpcAdapter) adapterLoop() {
   for {
      select {
          

      }
   }
}

func (a *RpcAdapter) Call(peerId string, req []byte, timeout int) ([]byte, error) {
   e := &RegistryEntry{}
   seq := a.seq.Next()
   e.req = &Message{seq:seq, data:req}
   e.done = make(chan int)
   a.rpcReg.put(e)
   a.GetTxChan() <- e
   if timeout == 0 {  //no timeout
      <-e.done
      return e.rsp.data, nil
   }
   to := make(chan int64, 1)
   t.timer.Schedule(timeout, to)
   select {
     case <-e.done:
        return e.rsp.data, nil
     case <-to:
        t.reqReg.del(seq)
        return "", ERR_RPC_TIMEOUT
   }
}


