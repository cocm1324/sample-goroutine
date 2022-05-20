package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cocm1324/sample-goroutine/pkg/pool"
)

func main() {
	// setting log
	logFile, err := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// starting
	log.Printf("main: starting up\n")
	ctx := context.Background()

	p := pool.New()
	poolCtx, poolCancel := context.WithCancel(ctx)
	p.Start(poolCtx)

	p.Create(poolCtx, 10)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("main: ctrl+c received, starting cleanup\n")
		poolCancel()

		time.Sleep(5 * time.Second)
		os.Exit(1)
	}()

	log.Printf("main: start up completed and running\n")
	for {
		time.Sleep(10 * time.Second)
	}
}
