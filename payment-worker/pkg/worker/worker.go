package workers

import (
	"context"
	"sync"
)

type Worker struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWorker() *Worker {
	ctx, cancel := context.WithCancel(context.Background())

	return &Worker{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *Worker) Run(job func(workerCtx context.Context)) {
	w.wg.Add(1)

	go func() {
		defer w.wg.Done()
		job(w.ctx)
	}()
}

func (w *Worker) Join() {
	w.wg.Wait()
}

func (w *Worker) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
}
