package database

const (
	workerCount = 64
)

var (
	jobQueue = make(chan *Job, workerCount)
)

func init() {
	for i := uint16(0); i < workerCount; i++ {
		go worker(i)
	}
}

func worker(id uint16) {
	//log.Debug().Uint16("Worker ID", id).Msg("database worker started")
	//defer log.Debug().Uint16("Worker ID", id).Msg("database worker stopped")

	for job := range jobQueue {
		job.Result <- job.Run()
		job.Done()
	}
}
