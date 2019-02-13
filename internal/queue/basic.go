package queue

import (
	"errors"
	"github.com/jjmschofield/gocrawl/internal/crawl"
)

type BasicQueue struct {
	channels Channels
	counters Counters
}

func NewBasicQueue() (queue *BasicQueue, err error) {
	return &BasicQueue{}, nil
}

func (q *BasicQueue) Start(worker crawl.QueueWorker, workerCount int) (results *chan crawl.WorkerResult, err error) {
	q.channels = Channels{
		jobs:    make(chan crawl.WorkerJob),
		Results: make(chan crawl.WorkerResult),
	}

	for i := 0; i < workerCount; i++ {
		go worker.Start(q.channels, &q.counters.Queue, &q.counters.Work)
	}

	return &q.channels.Results, nil
}

func (q *BasicQueue) Stop() (err error) {
	close(q.channels.Results)
	close(q.channels.jobs)
	return nil
}

func (q *BasicQueue) Push(job crawl.WorkerJob) (err error) {
	if q.channels.jobs == nil || q.channels.Results == nil  {
		return errors.New("queues are not open for use")
	}

	go func() {
		q.counters.Queue.Add(1)
		q.channels.jobs <- job
	}()

	return nil
}

func (q *BasicQueue) Counters() *Counters {
	return &q.counters
}
