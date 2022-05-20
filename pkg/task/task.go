package task

import (
	"context"
	"log"
	"time"
)

const DURATION int = 3000

type ID struct{}

type Task struct {
}

func New() *Task {
	return &Task{}
}

func (t *Task) Start(ctx context.Context, done chan int) {
	go t.start(ctx, done)
}

func (t *Task) start(ctx context.Context, done chan int) {
	start := time.Now()
	id := ctx.Value(ID{}).(int)
	log.Printf("task: starting %d\n", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("task: cancel signal received, canceling %d\n", id)
			return
		default:
			if time.Since(start) > time.Duration(DURATION)*time.Millisecond {
				log.Printf("task: done %d\n", id)
				done <- id
				return
			}
		}
	}
}
