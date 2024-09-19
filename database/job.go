package database

import (
	"sync"
	"time"

	"github.com/oomph-ac/api/errors"
)

var jobPool = sync.Pool{
	New: func() any {
		return &Job{}
	},
}

type Job struct {
	Run    func() interface{}
	Result chan interface{}
}

func (j *Job) Done() {
	j.Run = nil
	close(j.Result)
	jobPool.Put(j)
}

// RunJob returns a new job from the job pool.
func RunJob(r func() interface{}) (interface{}, *errors.APIError) {
	job := jobPool.Get().(*Job)
	job.Run = r
	job.Result = make(chan interface{}, 1)
	timeoutAt := time.Now().Add(time.Second * 10)

	// Try to get the job into the job queue so that a worker can handle the request. If there are no workers
	// available to handle the request, return an API error.
	select {
	case jobQueue <- job:
		// OK - the job went into the queue and is being handled by a worker.
	case <-time.After(time.Until(timeoutAt)):
		// We are at capacity and cannot handle this request at the time being.
		job.Done()
		return nil, errors.New(
			errors.APINoCapacity,
			"database workers at capacity",
			nil,
		)
	}

	// Though this should never happen, if for some reason the result channel is closed unexpectedly, we
	// need to return an error here.
	v, ok := <-job.Result
	if !ok {
		return nil, errors.New(
			errors.APIUnexpectedValue,
			"job result channel closed early",
			nil,
		)
	}

	// Check if the job returned an API error.
	if err, ok := v.(*errors.APIError); ok {
		return nil, err
	}

	// Horray - the job was successful! We can return this value with no errors.\
	return v, nil
}
