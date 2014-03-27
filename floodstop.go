package main

import (
	"time"
)

type Floodstop struct {
	ask  chan chan bool
	stop chan struct{}
}

func NewFloodstop(reset time.Duration, countMax int) (fs *Floodstop) {
	fs = &Floodstop{
		ask:  make(chan chan bool),
		stop: make(chan struct{}),
	}

	ticker := time.NewTicker(reset)
	counter := 0

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-fs.stop:
				return
			case retCh := <-fs.ask:
				counter++
				retCh <- (counter < countMax)
			case <-ticker.C:
				counter = 0
			}
		}
	}()

	return
}

func (fs *Floodstop) Stop() {
	fs.stop <- struct{}{}
}

func (fs *Floodstop) Ask() bool {
	ch := make(chan bool)
	fs.ask <- ch
	return <-ch
}
