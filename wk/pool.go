package wk

import (
	"context"
	"sync"
)

// Pool of worker
type Pool struct {
	ctx          context.Context
	cancel       context.CancelFunc
	numberWorker int
	wg           sync.WaitGroup
	ch           chan *Task
}

// NewPool create new worker pool
func NewPool(ctx context.Context, numberWorker int) (p *Pool) {
	if numberWorker <= 0 {
		numberWorker = 1
	}

	if ctx == nil {
		ctx = context.Background()
	}

	p = &Pool{
		numberWorker: numberWorker,
		ch:           make(chan *Task, numberWorker),
	}
	p.ctx, p.cancel = context.WithCancel(ctx)

	return
}

// Start workers
func (p *Pool) Start() {
	p.wg.Add(p.numberWorker + 1)
	for i := 0; i < p.numberWorker; i++ {
		go p.worker()
	}
}

// Do a task
func (p *Pool) Do(t *Task) {
	if p.ch != nil && t != nil {
		select {
		case <-p.ctx.Done():
		case p.ch <- t:
		}
	}
}

// Stop worker. Wait all task done.
func (p *Pool) Stop() {
	p.cancel()

	// wait child workers
	p.wg.Wait()
}

func (p *Pool) worker() {
	defer p.wg.Done()

	var task *Task
	for {
		select {
		case <-p.ctx.Done():
			return
		case task = <-p.ch:
			if task != nil {
				task.Execute()
			}
		}
	}
}
