package queue

import (
	"errors"
	"github.com/jjmschofield/gocrawl/internal/crawl"
)

type ReddisQueue struct {
	channels Channels
	counters Counters
	// redis address
	// redis client
	// is running
}

func NewReddisQueue() (queue *ReddisQueue, err error) {
	// set redis address
	return &ReddisQueue{}, nil
}

func (q *ReddisQueue) Start(worker crawl.QueueWorker, workerCount int) (results *chan crawl.WorkerResult, err error) {
	// Connect to redis
	// save client

	q.channels = Channels{
		jobs:    make(chan crawl.WorkerJob),
		Results: make(chan crawl.WorkerResult),
	}

	for i := 0; i < workerCount; i++ {
		go worker.Start(q.channels, &q.counters.Queue, &q.counters.Work)
	}

	go q.pollForJobs()

	return &q.channels.Results, nil
}

func (q *ReddisQueue) Stop() (err error) {
	// set running false
	close(q.channels.Results)
	close(q.channels.jobs)
	return nil
}

func (q *ReddisQueue) Push(job crawl.WorkerJob) (err error) {
	// if running instead

	if q.channels.jobs == nil || q.channels.Results == nil  {
		return errors.New("queues are not open for use")
	}

	// q.counters.Queue.Add(1)
	// Push to reddis

	return nil
}

func (q *ReddisQueue) Counters() *Counters {
	return &q.counters
}

func (q *ReddisQueue) pollForJobs(){
	// While we are open
	// Poll redis
	// Pop job
	// Push jobs into jobs channel
	// Block when there are not enough workers
	// close client
}
