package queue

import (
	"errors"
	"github.com/jjmschofield/gocrawl/internal/counters"
)

type BasicQueue struct {
	channels Channels
	counters Counters
}

func NewBasicQueue() (queue *BasicQueue) {
	return &BasicQueue{
		counters: Counters{
			Work:  &counters.AtomicInt64{},
			Queue: &counters.AtomicInt64{},
		},
	}
}

func (q *BasicQueue) Start(worker QueueWorker, workerCount int) (results chan WorkerResult, err error) {
	q.channels = Channels{
		Jobs:    make(chan WorkerJob),
		Results: make(chan WorkerResult),
	}

	for i := 0; i < workerCount; i++ {
		go worker.Start(q.channels, q.counters.Queue, q.counters.Work)
	}

	return q.channels.Results, nil
}

func (q *BasicQueue) Stop() {
	close(q.channels.Results)
	close(q.channels.Jobs)
}

func (q *BasicQueue) Push(job WorkerJob) (err error) {
	if q.channels.Jobs == nil || q.channels.Results == nil {
		return errors.New("cannot push to queue, channels are not open for use")
	}

	go func() {
		q.counters.Queue.Add(1)
		q.channels.Jobs <- job
	}()

	return nil
}

func (q *BasicQueue) Counters() Counters {
	return q.counters
}
