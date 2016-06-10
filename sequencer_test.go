package hrpc

import (
   "testing"
)

func TestSeq32Next(t *testing.T) {
   s := NewSequencer32(int32(0))
   s1 := s.Next()
   if s1 != int32(1) {
      t.Fatal("Next() return wrong value")
   }
}

func TestSeq64Next(t *testing.T) {
   s := NewSequencer64(int64(0))
   s1 := s.Next()
   if s1 != int64(1) {
      t.Fatal("Next() return wrong value")
   }
}
