package clients

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/streamdp/ccd/domain"
)

// Task does all the data mining run
type Task struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Interval int64  `json:"interval"`
	done     chan struct{}
}
type Tasks map[string]*Task

func (t *Task) run(r RestClient, l *log.Logger, dataPipe chan *domain.Data) {
	timer := time.NewTimer(time.Duration(atomic.LoadInt64(&t.Interval)) * time.Second)
	go func() {
		defer close(t.done)
		for {
			timer.Reset(time.Duration(atomic.LoadInt64(&t.Interval)) * time.Second)
			select {
			case <-t.done:
				return
			case <-timer.C:
				data, err := r.Get(t.From, t.To)
				if err != nil {
					l.Println(err)
					continue
				}
				dataPipe <- data
			}
		}
	}()
}

func (t *Task) close() {
	t.done <- struct{}{}
}
