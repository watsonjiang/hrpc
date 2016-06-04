package hrpc


//interface to both tcp transporter and udp transporter
type Trans interface{
   AddPeer(peerInfo string) error

   Send(m *Message)
  
   RegisteListener(l TransListener)

   Start()

   Close()
}

type TransListener interface {
   OnReqArrival(m *Message)
   OnRspArrival(m *Message)
}

type TransConfig struct {
   LocalPeerInfo  string  //{"id":"qs01", "ct":"61.22.34.14:9900", "cnc":"10.22.34.53:3340"}
   MaxConns int      //valid for tcp trans, max number of tcp connection
   MaxSendRate int   //in k-bytes per second, 0 - disabled
   BigMsgSize int    //message with size>BigMsgSize will be transfered
                     //by sepqrate connection
}

func NewTransDefaultCfg() *TransConfig {
   return &TransConfig{}
}

func NewTcpTrans(cfg *TransConfig) Trans {
   return &TcpTransImpl{}
}
