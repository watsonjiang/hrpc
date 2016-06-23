package hrpc

import (
   "net"
   "time"
   "bytes"
   "io"
)

type BaseMockConn struct{
   Tx bytes.Buffer
   Rx bytes.Buffer
}

func (c *BaseMockConn) getTx() []byte{
   return c.Tx.Bytes()
}

func (c *BaseMockConn) Read(b []byte) (n int, err error) {
   return c.Rx.Read(b)
}

func (c *BaseMockConn) Write(b []byte) (n int, err error) {
   return c.Tx.Write(b)
}

func (c *BaseMockConn) Close() error {
   return nil
}

func (c *BaseMockConn) LocalAddr() net.Addr {
   return nil
}

func (c *BaseMockConn) RemoteAddr() net.Addr {
   return nil
}

func (c *BaseMockConn) SetDeadline(t time.Time) error {
   return nil
}

func (c *BaseMockConn) SetReadDeadline(t time.Time) error {
   return nil
}

func (c *BaseMockConn) SetWriteDeadline(t time.Time) error {
   return nil
}

//a conn obj act as a normal conn.
//no err for both read and write
type MockNormalConn struct{
   BaseMockConn
}

//a conn obj any write operation will return timeout error
type MockWriteTimeoutConn struct {
   BaseMockConn
}

type TimeoutError struct {
}

func (e *TimeoutError) Error() string {
   return "Timeout error"
}

func (e *TimeoutError) Timeout() bool {
   return true
}

func (e *TimeoutError) Temporary() bool {
   return false
}

func (c *MockWriteTimeoutConn) Write(b []byte) (n int, err error) {
   return 0, &TimeoutError{}
}

type MockReadEofConn struct {
   BaseMockConn
}

func (c *MockReadEofConn) Read(b []byte) (n int, err error) {
   return 0, io.EOF
}
