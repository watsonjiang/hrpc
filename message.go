package hrpc

const (
   MSG_TYPE_REQ
   MSG_TYPE_RSP
   )

//message
type Message struct {
   seq int32
   mtype int
   data []byte
}

func (r *Message) isReqMsg() bool {
   return r.mtype == MSG_TYPE_REQ 
}

func (r *Message) Bytes() []byte {
   buf := new(bytes.Buffer)
   binary.Write(buf, binary.LittleEndian, r.seq)
   binary.Write(buf, binary.LittleEndian, int16(len(r.data)))
   binary.Write(buf, binary.LittleEndian, r.data)
   return buf.Bytes()
}

func decodeMessage(buf []byte) *Message{
   r := &Request{}
   rd := bytes.NewReader(buf)
   binary.Read(rd, binary.LittleEndian, &r.seq)
   var ld int16
   binary.Read(rd, binary.LittleEndian, &ld)
   r.data = make([]byte, ld)
   binary.Read(rd, binary.LittleEndian, &r.data)
   return r
}
