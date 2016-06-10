package hrpc

import (
   "sync/atomic"
)

type Sequencer32 struct {
   seq int32
}

func NewSequencer32(init int32) *Sequencer32{
   return &Sequencer32{seq:init}
}

func (s *Sequencer32) Next() int32{
   return atomic.AddInt32(&s.seq, 1)
}

type Sequencer64 struct {
   seq int64
}

func NewSequencer64(init int64) *Sequencer64{
   return &Sequencer64{seq:init}
}

func (s *Sequencer64) Next() int64{
   return atomic.AddInt64(&s.seq, 1)
}
