package hrpc

import (
   "testing"
)

func client() {
   t := NewTrans(NewTransDefaultCfg())
   t.Dial("127.0.0.1:9999")
   data := []byte{"hello"}
   for i:=0;i<3;i++{
      go func(){
         t.Call(data, 0)
      }()
   }
}

func handler(req []byte) []byte {
   return []byte{"world"}
}

func server() {
   t := NewTrans(NewTransDefaultCfg())
   t.registerRpcHandler(handler)
   t.Listen("127.0.0.1:9999")
   for i:=0; i<3;i++{
      go func() {
         for {
            t.PumpRequest()
         }
      }()
   }
}

func TestExample(t *testing.T) {
   go server()
   go client()
}
