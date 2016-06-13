package hrpc

import (
   "testing"
   "encoding/hex"
   "bytes"
)

func TestRequest(t *testing.T) {
   req := NewRequest()
   if !req.IsReqMsg() {
      t.Error("IsReqMsg not correct")
   }
   data := []byte("1234567890")
   req.data = data
   req.seq = 1
   var raw bytes.Buffer
   cnt, err := req.WriteTo(&raw)
   t.Logf("req: %s raw: %s", req, hex.EncodeToString(raw.Bytes()))
   if cnt!=raw.Len() || err!=nil {
      t.Error("WriteTo not correct.", cnt, raw.Len())
   }
   m := NewMessage()
   if err:=m.ReadFrom(&raw);err!=nil {
      t.Fatal("ReadFrom return error.", err)
   }
   if m.seq != 1 {
      t.Error("ReadFrom not correct.")
   }
   if !m.IsReqMsg() {
      t.Error("ReadFrom not correct.")
   }
   if !bytes.Equal(data, m.data) {
      t.Errorf("ReadFrom not correct. exp:1234567890 got:%s",
              string(m.data))
   }
}
