package hrpc

import (
   "fmt"
   "bytes"
   "encoding/binary"
   "encoding/hex"
   "io"
)

const (
   MSG_RSP_BIT = 0x01
)

//message
type Message struct {
   seq int32
   mtype int8
   peerId string
   data bytes.Buffer
}

func NewMessage() *Message {
   return &Message{}
}

func (m *Message) String() string {
   return fmt.Sprintf("seq %v data %s", m.seq,
                      hex.EncodeToString(m.data.Bytes()))
}

func (m *Message) IsReqMsg() bool {
   if 0 == m.mtype & MSG_RSP_BIT {
      return false
   }
   return true
}

func (m *Message) WriteTo(w io.Writer) (int, error) {
   if err:=binary.Write(w, binary.LittleEndian, m.seq); err!=nil {
      return 0, err
   }
   if err:=binary.Write(w, binary.LittleEndian, m.mtype); err!=nil {
      return 0, err
   }
   if err:=binary.Write(w, binary.LittleEndian, uint16(m.data.Len())); err!=nil {
      return 0, err
   }
   if _, err:=m.data.WriteTo(w);err!=nil {
      return 0, err
   }
   cnt := 4 + 1 + 2 + m.data.Len()
   return cnt, nil
}

func (m *Message) ReadFrom(r io.Reader) error {
   if err:=binary.Read(r, binary.LittleEndian, &m.seq);err!=nil {
      return err
   }
   if err:=binary.Read(r, binary.LittleEndian, &m.mtype);err!=nil{
      return err
   }
   var ld uint16
   if err:=binary.Read(r, binary.LittleEndian, &ld);err!=nil{
      return err
   }
   if uint16(m.data.Cap()) < ld {
      m.data.Grow(int(ld))
   }
   m.data.Reset()
   if _, err:=io.CopyN(&m.data, r, int64(ld));err!=nil {
      return err
   }
   return nil
}

func (m *Message) MakeResponse() *Message {
   rsp := &Message{}
   rsp.seq = m.seq
   rsp.mtype &= MSG_RSP_BIT
   return rsp
}

func NewRequest() *Message{
   m := &Message{}
   m.mtype &= ^MSG_RSP_BIT
   return m
}

