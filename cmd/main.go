package main

import (
	"context"
	"github.com/cingozr/yemek_sepeti/handlers"
	"github.com/cingozr/yemek_sepeti/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	fromTimerToJsonChnModel chan *struct{}
	timerRecordService      *service.TimerRecording
	JsonStorageService      *service.JsonStorage
	timerDuration           = 60

	sigCh  chan os.Signal
	ctx    context.Context
	cancel context.CancelFunc
)

func init() {
	sigCh = make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	ctx, cancel = context.WithCancel(context.Background())

	fromTimerToJsonChnModel = make(chan *struct{})
	timerRecordService = service.NewTimerRecording(ctx, fromTimerToJsonChnModel, timerDuration)
	go timerRecordService.SetFileInInterval()

	JsonStorageService = service.NewJsonStorage(ctx, fromTimerToJsonChnModel)

	go JsonStorageService.IntervalSave()
}

func main() {
	defer cancel()


	http.HandleFunc("/set_key", handlers.SaveKey)
	http.HandleFunc("/get_key", handlers.GetKey)
	http.HandleFunc("/flush_memory", handlers.FlushMemory)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		cancel()
		return
	}

	select {
	case <-ctx.Done():
	case <-sigCh:
		cancel()
	}
}
