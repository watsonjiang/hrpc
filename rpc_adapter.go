package hrpc

import (
   "sync"
   "errors"
)

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

func NewRpcRegistry() *RpcRegistry {
   return &RpcRegistry{
      reg:make(map[int32] *RpcRegistryEntry),
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
      delete(r.reg, seq)
      return e
   }
   return nil
}

type RpcAdapter struct {
   trans Trans
   seq *Sequencer32
   rpcReg *RpcRegistry
   rpcHandler RpcHandlerFunc
   timer *TimerQueue
}

func (a *RpcAdapter) RegisterRpcHandler(f RpcHandlerFunc) {
   a.rpcHandler = f
}

func (a *RpcAdapter) OnReqArrival(req *Message) {
   data := a.rpcHandler(req.data)
   rsp := req.MakeResponse()
   rsp.data = data
   a.trans.Send(rsp)
}

func (a *RpcAdapter) OnRspArrival(rsp *Message) {
   if e:=a.rpcReg.del(rsp.seq);e!=nil {
      e.rsp = rsp
      e.done <- 1
   }
}

func (a *RpcAdapter) Call(peerId string, req []byte, timeout int) ([]byte, error) {
   e := &RpcRegistryEntry{}
   seq := a.seq.Next()
   e.req = &Message{seq:seq, data:req}
   e.done = make(chan int)
   a.rpcReg.add(seq, e)
   a.trans.Send(e.req)
   if timeout == 0 {  //no timeout
      <-e.done
      return e.rsp.data, nil
   }
   to := make(chan int64, 1)
   a.timer.Schedule(int64(timeout), to)
   select {
     case <-e.done:
        return e.rsp.data, nil
     case <-to:
        a.rpcReg.del(seq)
        return nil, ERR_RPC_TIMEOUT
   }
}


