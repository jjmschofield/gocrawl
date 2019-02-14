package queue

import (
	"errors"
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

func (q *ReddisQueue) Start(worker QueueWorker, workerCount int) (results *chan WorkerResult, err error) {
	// Connect to redis
	// save client

	q.channels = Channels{
		Jobs:    make(chan WorkerJob),
		Results: make(chan WorkerResult),
	}

	for i := 0; i < workerCount; i++ {
		go worker.Start(q.channels, q.counters.Queue, q.counters.Work)
	}

	go q.pollForJobs()

	return &q.channels.Results, nil
}

func (q *ReddisQueue) Stop(){
	// set running false
	close(q.channels.Results)
	close(q.channels.Jobs)
}

func (q *ReddisQueue) Push(job WorkerJob) (err error) {
	// if running instead

	if q.channels.Jobs == nil || q.channels.Results == nil  {
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
	// Push Jobs into Jobs channel
	// Block when there are not enough workers
	// close client
}
