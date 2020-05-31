package wk

import (
	"context"
)

// Task represents a task
type Task struct {
	ctx      context.Context
	info     interface{}
	executor func(context.Context, interface{}) error
}

// NewTask create new task
func NewTask(ctx context.Context, taskInfo interface{}, executor func(context.Context, interface{}) error) *Task {
	return &Task{
		ctx:      ctx,
		info:     taskInfo,
		executor: executor,
	}
}

// Execute task
func (t *Task) Execute() {
	if t.executor != nil {
		_ = t.executor(t.ctx, t.info)
	}
}
