package hrpc

import (
   "testing"
   "encoding/hex"
)

func TestSendHandshake(t *testing.T) {
   info := `{"Id":"xxxxx", "Addr":["127.0.0.1:1234"]}`
   l_peer := NewPeer(info)
   c := &MockNormalConn{}
   err := sendHandshake(l_peer, c)
   if err!=nil {
      t.Error("sendHandshake return non-nil")
   }
   t.Log("tx", hex.EncodeToString(c.getTx()))
}
