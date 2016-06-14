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
   if cnt!=int64(raw.Len()) || err!=nil {
      t.Error("WriteTo not correct.", cnt, raw.Len())
   }
   m := NewMessage()
   if n,err:=m.ReadFrom(&raw);err!=nil {
      t.Fatal("ReadFrom return error.", err)
   }else if n != cnt {
      t.Error("ReadFrom return cnt not correct. exp:", cnt, "act:", n)
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


