package hrpc

import (
   "testing"
   "bytes"
)

func TestNewPeer(t *testing.T) {
   info := `{"Id":"xxxxx", "Addr":["127.0.0.1:1234"]}`
   p := NewPeer(info)
   if p.Id != "xxxxx" {
      t.Error("NewPeer Id not correct.")
   }
   if len(p.Addr) != 1 && p.Addr[0] != "127.0.0.1:1234" {
      t.Error("NewPeer Addr not correct.")
   }
}

type _testListener struct {
   _added *Peer
   _update_old, _update_new *Peer
}

func (l *_testListener) OnPeerAdded(p *Peer) {
   l._added = p
}

func (l *_testListener) OnPeerUpdated(oldv, newv *Peer) {
   l._update_old, l._update_new = oldv, newv
}

func TestRegListener(t *testing.T){
   l := &_testListener{}
   reg := NewPeerRegistry(l)
   info := `{"Id":"xxxxx", "Addr":["127.0.0.1:1234"]}`
   p := NewPeer(info)
   reg.Put(p)
   if l._added != p {
      t.Error("OnPeerAdded not correct.")
   }
   reg.Put(p)
   if l._update_old != p && l._update_new != p {
      t.Error("OnPeerUpdated not correct.")
   }
}

func TestRegPutGet(t *testing.T){
   reg := NewPeerRegistry(nil)
   info := `{"Id":"xxxxx", "Addr":["127.0.0.1:1234"]}`
   p := NewPeer(info)
   reg.Put(p)
   r := reg.Get("xxxxx")
   if len(r.Addr) != 1 && r.Addr[0] != "127.0.0.1:1234" {
      t.Error("Get not correct.")
   }
}

func TestPeerMarshalUnmarshal(t *testing.T) {
   peerInfo := `{"Id":"peer1", "Addr":["127.0.0.1:1234"]}`
   peer := NewPeer(peerInfo)
   if peer == nil {
      t.Fatal("NewPeer return nil")
   }
   var raw = new(bytes.Buffer)
   n, err := peer.WriteTo(raw)
   if err!=nil || n != int64(raw.Len()) {
      t.Error("peer WriteTo not correct.")
   }
   peer1 := new(Peer)
   l := int64(raw.Len())
   n, err = peer1.ReadFrom(raw)
   if err!=nil || n != l {
      t.Error("peer ReadFrom not correct.", err, n, l)
   }
   t.Logf("%v", peer1)
}

