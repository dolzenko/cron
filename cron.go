package cron

import (
	"sync"
	"time"
)

type cbF func() error
type errCbF func(error)

type task struct {
	dur time.Duration
	cb  cbF
}

type cron struct {
	tasks chan task
	errCb errCbF
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
func Every(dur time.Duration, cbs ...cbF) *cron {
	return defaultCron.Every(dur, cbs...)
}

// OnError registers task error handler
func OnError(cb errCbF) *cron {
	return defaultCron.OnError(cb)
}

// Every schedules execution of callback(s)
func (c *cron) Every(dur time.Duration, cbs ...cbF) *cron {
	for _, cb := range cbs {
		c.tasks <- task{dur, cb}
	}
	return c
}

// OnError registers task error handler
func (c *cron) OnError(cb errCbF) *cron {
	c.errCb = cb
	return c
}

func (c *cron) loop() {
	for t := range c.tasks {
		go c.taskLoop(t)
	}
}

func (c *cron) taskLoop(t task) {
	t.run(c.errCb)
	for range time.NewTicker(t.dur).C {
		t.run(c.errCb)
	}
}

func (t task) run(errCb errCbF) {
	err := t.cb()
	if err != nil && errCb != nil {
		errCb(err)
	}
}
