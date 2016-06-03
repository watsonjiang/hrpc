package hrpc

import (
   "atomic"
   )

type Sequencer struct {
   seq int32
}

func NewSequencer(init int32) *Sequencer{
   return &Sequencer{seq:init}
}

func (s *Sequencer) Next() int32{
   return atomic.AddInt32(&s.seq, 1)
}
