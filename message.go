package hrpc

import (
   "fmt"
   "encoding/binary"
   "encoding/hex"
   "io"
   "strconv"
)

const (
   MSG_RSP_BIT = 0x01
   MSG_HANDSHAKE_BIT = 0x02
)

//message
type Message struct {
   seq int32
   mtype int8
   peerId string
   data []byte
}

func NewMessage() *Message {
   return &Message{data:nil}
}

func (m *Message) String() string {
   return fmt.Sprintf("{seq:%v, mtype:%08s, data:%s}", m.seq,
                      strconv.FormatInt(int64(m.mtype), 2),
                      hex.EncodeToString(m.data))
}

func (m *Message) IsReqMsg() bool {
   if 0 == m.mtype & MSG_RSP_BIT {
      return true
   }
   return false
}

func (m *Message) IsHandshakeMsg() bool {
   if 0 == m.mtype & MSG_HANDSHAKE_BIT {
      return false
   }
   return true
}

func (m *Message) WriteTo(w io.Writer) (int64, error) {
   cnt := int64(0)
   if err:=binary.Write(w, binary.LittleEndian, m.seq); err!=nil {
      return cnt, err
   }
   cnt += int64(4)
   if err:=binary.Write(w, binary.LittleEndian, m.mtype); err!=nil {
      return cnt, err
   }
   cnt += int64(1)
   if err:=binary.Write(w, binary.LittleEndian, uint16(len(m.data))); err!=nil {
      return cnt, err
   }
   cnt += int64(2)
   if _, err:=w.Write(m.data);err!=nil {
      return cnt, err
   }
   cnt += int64(len(m.data))
   return cnt, nil
}

func (m *Message) ReadFrom(r io.Reader) (int64, error) {
   cnt := int64(0)
   if err:=binary.Read(r, binary.LittleEndian, &m.seq);err!=nil {
      return cnt, err
   }
   cnt += int64(4)
   if err:=binary.Read(r, binary.LittleEndian, &m.mtype);err!=nil{
      return cnt, err
   }
   cnt += int64(1)
   var ld uint16
   if err:=binary.Read(r, binary.LittleEndian, &ld);err!=nil{
      return cnt, err
   }
   cnt += int64(2)
   if uint16(cap(m.data)) < ld {
      m.data = make([]byte, ld)
   }else{
      m.data = m.data[:ld]
   }
   if _, err:=io.ReadFull(r, m.data);err!=nil {
      return cnt, err
   }
   cnt += int64(len(m.data))
   return cnt, nil
}

func (m *Message) MakeResponse() *Message {
   rsp := &Message{}
   rsp.seq = m.seq
   rsp.mtype &= MSG_RSP_BIT
   return rsp
}

func NewRequest() *Message{
   m := NewMessage()
   m.mtype &= ^MSG_RSP_BIT
   return m
}

