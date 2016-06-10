package hrpc

import (
   "testing"
)

func TestSeqNext(t *testing.T) {
   s := NewSequencer(int32(0))
   s1 := s.Next()
   if s1 != int32(1) {
      t.Fatal("Next() return wrong value")
   }
}
