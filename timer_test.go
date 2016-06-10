package hrpc

import (
   "testing"
   "time"
)

func TestTimerStartStop(t *testing.T){
   tq := NewTimerQueue()
   tq.Start()
   t.Log("TimerQueue started.")
   tq.Stop()
   t.Log("TimerQueue stoped.")
}

func TestTimerSchedule(t *testing.T) {
   tq := NewTimerQueue()
   tq.Start()
   ch := make(chan int64, 1)
   tq.Schedule(1000, ch)
   select {
      case <-time.After(3 * time.Second):
         t.Fatal("TimerQueue schedule not work.")
      case <-ch:
   }
   tq.Stop()
}
