package queue

import (
	"github.com/jjmschofield/gocrawl/internal/counters"
	"github.com/jjmschofield/gocrawl/internal/crawl"
)

type Queue interface {
	Start(worker crawl.QueueWorker, workerCount int) (results *chan crawl.WorkerResult, err error)
	Stop() (err error)
	Push(job crawl.WorkerJob) (err error)
	Counters() *Counters
}

type Channels struct {
	jobs chan crawl.WorkerJob
	Results chan crawl.WorkerResult
}

type Counters struct{
	Queue counters.AtomicInt64
	Work counters.AtomicInt64
}