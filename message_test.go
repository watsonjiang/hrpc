package hrpc

import (
   "testing"
   "encoding/hex"
)

func TestRequest(t *testing.T) {
   req = NewRequest()
   if !req.IsReqMsg() {
      t.Error("IsReqMsg not correct")
   }
   data = []byte("1234567890")
   req.data = data
   raw = req.Bytes()
   t.Logf("req: %s %s", req, hex.EncodingToString(raw))
   req1 = decodeMessage(raw) 

}
