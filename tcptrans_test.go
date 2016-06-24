package hrpc

import (
   "testing"
   "encoding/hex"
   "net"
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

//send return with timeout error
func TestSendHandshake1(t *testing.T) {
   info := `{"Id":"xxxxx", "Addr":["127.0.0.1:1234"]}`
   l_peer := NewPeer(info)
   c := &MockWriteTimeoutConn{}
   err := sendHandshake(l_peer, c)
   if err==nil{
      t.Error("sendHandshake not return nil")
   }
   if err1, ok:=err.(net.Error);ok {
      if err1.Timeout() {
         return
      }
   }
   t.Error("Is not timeout error.")
}

func TestReadHandshake(t *testing.T) {
   c := &MockNormalConn{}
   tmp, _ := hex.DecodeString("00000000021600057878787878010e3132372e302e302e313a31323334")
   c.Rx.Write(tmp)
   if p, err := readHandshake(c);err!=nil {
      t.Error("readHandshake fail.", err)
   }else{
      t.Log(p)
      if p.Id != "xxxxx" {
         t.Error("peer Id not correct.")
      }
      if len(p.Addr)==1 && p.Addr[0]=="127.0.0.1:1234" {
         t.Log("peer is ok")
      }else{
         t.Error("peer addr not correct.")
      }
   }
}

//read EOF
func TestReadHandshake1(t *testing.T) {
   c := &MockReadEofConn{}
   if _, err := readHandshake(c);err==nil{
      t.Error("readHandshake not return err.")
   }else{
      t.Log(err)
   }
}

//read invalid handshake message
func TestReadHandshake2(t *testing.T) {
   c := &MockNormalConn{}
   tmp, _ := hex.DecodeString("00000000001600057878787878010e3132372e302e302e313a31323334")
   c.Rx.Write(tmp)
   if _, err := readHandshake(c);err==nil{
      t.Error("readHandshake not return err.")
   }else{
      t.Log(err)
   }
}

