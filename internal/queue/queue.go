package queue

import (
	"github.com/jjmschofield/gocrawl/internal/counters"
)

type Queue interface {
	Start(worker QueueWorker, workerCount int) (results chan WorkerResult, err error)
	Stop()
	Push(job WorkerJob) (err error)
	Counters() Counters
}

type Channels struct {
	Jobs    chan WorkerJob
	Results chan WorkerResult
}

type Counters struct {
	Queue *counters.AtomicInt64
	Work  *counters.AtomicInt64
}
