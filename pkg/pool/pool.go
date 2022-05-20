package pool

import (
	"context"
	"log"
	"time"

	"github.com/cocm1324/sample-goroutine/pkg/task"
)

const DONE_BUFFER = 5
const CREATE_DURATION = 1500

type Pool struct {
	lastIssued int
	done       chan int
}

func New() *Pool {
	return &Pool{lastIssued: 0}
}

func (p *Pool) Start(ctx context.Context) {
	p.done = make(chan int, DONE_BUFFER)
	go p.start(ctx)
}

func (p *Pool) start(ctx context.Context) {
	log.Printf("pool: starting up\n")
L:
	for {
		select {
		case <-ctx.Done():
			log.Printf("pool: cancel signal received\n")
			break L
		case id := <-p.done:
			log.Printf("pool: task done signal received %d\n", id)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	log.Printf("pool: pool gracefully shut down\n")
}

func (p *Pool) Create(ctx context.Context, count int) {
	go p.createWait(ctx, count)
}

func (p *Pool) createWait(ctx context.Context, count int) {
	log.Printf("create: starting to create %d task(s)\n", count)
	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			log.Printf("create: cancel signal received\n")
			return
		default:
			p.lastIssued += 1
			t := task.New()
			taskContext := context.WithValue(ctx, task.ID{}, p.lastIssued)
			t.Start(taskContext, p.done)
			log.Printf("create: creating task %d\n", p.lastIssued)
			time.Sleep(CREATE_DURATION * time.Millisecond)
		}
	}
	log.Printf("create: completed to create %d tasks(s)\n", count)
}
