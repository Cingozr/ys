package service

import (
	"context"
	"sync"
	"time"
)

type TimerRecording struct {
	wg                 sync.WaitGroup
	ctx                context.Context
	duration           int
	fromTimerToJsonChn chan *struct{}
}

func NewTimerRecording(ctx context.Context, fromTimerToJsonChn chan *struct{}, duration int) *TimerRecording {
	return &TimerRecording{
		ctx:                ctx,
		fromTimerToJsonChn: fromTimerToJsonChn,
		duration: duration,
	}
}

func (tr *TimerRecording) SetFileInInterval() {
	tr.wg.Add(1)
	go func() {
		for {
			select {
			case <-time.NewTimer(time.Duration(tr.duration) * time.Second).C:
				tr.fromTimerToJsonChn <- &struct{}{}
			case <-tr.ctx.Done():
				tr.wg.Done()
				break
			}
		}
	}()
	tr.wg.Wait()
}
