package scheduler

import "time"

var Instance = New()

type Scheduler struct {
}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Schedule(every time.Duration, do func()) {
	go func() {
		var ticker = time.Tick(every)
		for {
			do()
			<-ticker
		}
	}()
}
