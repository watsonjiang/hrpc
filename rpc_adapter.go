package hrpc

type RpcHandlerFunc func(req []byte) []byte

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
   rpcHandler RpcHandlerFunc
}

func (a *RpcAdapter) RegisterRpcHandler(f RpcHandlerFunc) {
   a.rpcHandler = f
}

func (a *RpcAdapter) OnReqArrival(m *Message) *Message {
   data = a.rpcHandler(m.data)
   return &Message{seq:m.seq, data:data}
}

func (a *RpcAdapter) OnRspArrival(m *Message) {
   if e:=a.rpcReg.del(m.seq);e!=nil {
      e.rsp = m  
      e.done <- 1
   }
}

func (a *RpcAdapter) Call(peerId string, req []byte, timeout int) ([]byte, error) {
   e := &RegistryEntry{}
   seq := a.seq.Next()
   e.req = &Message{seq:seq, data:req}
   e.done = make(chan int)
   a.rpcReg.add(e)
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
        t.rpcReg.del(seq)
        return nil, ERR_RPC_TIMEOUT
   }
}


