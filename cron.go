package cron

import (
	"sync"
	"time"
)

type cbF func() error

type task struct {
	dur time.Duration
	cb  cbF
}

type cron struct {
	tasks chan task
	errCb func(error)
	sync.Mutex
}

// NewCron returns ready to use (and running) scheduler
func NewCron() *cron {
	c := &cron{
		tasks: make(chan task),
	}
	go c.loop()
	return c
}

var defaultCron = NewCron()

// Every schedules execution of callback(s)
func Every(dur time.Duration, cbs ...cbF) {
	defaultCron.Every(dur, cbs...)
}

// OnError registers task error handler
func OnError(cb func(error)) {
	defaultCron.OnError(cb)
}

// Every schedules execution of callback(s)
func (c *cron) Every(dur time.Duration, cbs ...cbF) {
	for _, cb := range cbs {
		c.tasks <- task{dur, cb}
	}
}

// OnError registers task error handler
func (c *cron) OnError(cb func(error)) {
	c.errCb = cb
}

func (c *cron) loop() {
	for t := range c.tasks {
		go c.run(t)
	}
}

func (c *cron) run(t task) {
	for range time.NewTicker(t.dur).C {
		if err := t.cb(); err != nil && c.errCb != nil {
			c.errCb(err)
		}
	}
}
